from __future__ import unicode_literals
import sys
import os
import argparse
import logging
import datetime
import csv
import json
import time
import math
import requests as http_requests

parentdir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
os.sys.path.insert(0, parentdir)

SPREAD = 0.01
# Test Server
# SERVER_URL = "http://localhost:8080/api/"
# Deployed Server
SERVER_URL = 'http://cloud-vm-45-112.doc.ic.ac.uk:8080/api/'
REQUEST_TYPES = {
    'order': 'order',
    'price': 'highest-bid-lowest-ask',
    'cancel': 'cancel-order'
}

def create_order_json(mkt_mkr_id, ticker, amount, ord_type, limit_price):
    return {
        'userId': mkt_mkr_id,
        'equityTicker': ticker,
        'amount': amount,
        'orderType': ord_type,
        'limitPrice': limit_price
    }


def create_price_json(ticker):
    return {
        'ticker': ticker
    }


def create_cancel_order_json(limit_price, mkt_mkr_id, ticker, bid):
    return {
        'limitPrice': limit_price,
        'userId': mkt_mkr_id,
        'ticker': ticker,
        'bid': bid
    }

def generate_with_depth(prices):
    return [(prices[i], 10**(i+1)) for i in range(len(prices))]


def populate_bids(highest_bid):
    prices = [highest_bid-(highest_bid*x/1000.0) for x in range(0,5)]
    return generate_with_depth(prices)


def populate_asks(lowest_ask):
    prices = [lowest_ask+(lowest_ask*x/1000.0) for x in range(0,5)]
    return generate_with_depth(prices)


def fill_missing_bids(bid_target, bid_now):
    steps = int(math.ceil((bid_target - bid_now)/(bid_target/1000.0)))
    prices = [bid_target-(bid_target*x/1000.0) for x in range(0,steps)]
    return generate_with_depth(prices)


def fill_missing_asks(ask_target, ask_now):
    steps = int(math.ceil((ask_now - ask_target)/(ask_target/1000.0)))
    prices = [ask_target+(ask_target*x/1000.0) for x in range(0,steps)]
    return generate_with_depth(prices)


def send_requests(requests, req_type):
    responses = []
    for request in requests:
        request_data = request
        print request_data
        response = http_requests.post(SERVER_URL + REQUEST_TYPES[req_type], json=request_data)
        print response.content
        try:
            responses.append(json.loads(response.content))
        except:
            print "no json response"

    return responses


def get_highest_bid(target_price):
    return target_price*(1-SPREAD/2.0)


def get_lowest_ask(target_price):
    return target_price*(1+SPREAD/2.0)


def sleep_and_get_prices(tickers):
    time.sleep(15)
    prices = send_requests([create_price_json(ticker) for ticker in tickers], 'price')

    # if len(tickers) != len(prices):
        # logging.error('Tickers and prices lengths do not match')

    print prices
    tickers_bids_asks = { tickers[i]: (prices[i]['highestBid'], prices[i]['lowestAsk']) for i in range(len(prices))}
    print tickers_bids_asks
    return tickers_bids_asks

def should_intervene(highest_bid, lowest_ask):
    diff = lowest_ask - highest_bid
    midpoint = highest_bid + diff/2.0

    return (diff/midpoint) >= SPREAD

def derive_price(highest_bid, lowest_ask):
    return lowest_ask - (lowest_ask - highest_bid)/2.0

def init(tickers_prices, market_maker_id):
    requests = []
    for ticker, target_price in tickers_prices.iteritems():
        bids = populate_bids(get_highest_bid(target_price))
        asks = populate_asks(get_lowest_ask(target_price))

        for bid_price, bid_amount in bids:
            requests.append(create_order_json(market_maker_id, ticker, bid_amount, "limitBid", bid_price))

        for ask_price, ask_amount in asks:
            requests.append(create_order_json(market_maker_id, ticker, ask_amount, "limitAsk", ask_price))

    send_requests(requests, 'order')


def intervene(market_maker_id, ticker, price_data):
    hb_now, la_now, hb_target, la_target = price_data
    requests = []
    if hb_now < hb_target:
        bids = fill_missing_bids(hb_target, hb_now)
        for bid_price, bid_amount in bids:
            requests.append(create_order_json(market_maker_id, ticker, bid_amount, "limitBid", bid_price))

    if la_now > la_target:
        asks = fill_missing_asks(la_target, la_now)
        for ask_price, ask_amount in asks:
            requests.append(create_order_json(market_maker_id, ticker, ask_amount, "limitAsk", ask_price))

    send_requests(requests, 'order')


def start_market_makers(tickers_prices, market_maker_ids):
    for market_maker_id in market_maker_ids:
        init(tickers_prices, market_maker_id)

    while True:
        for market_maker_id in market_maker_ids:
            tickers_bids_asks = sleep_and_get_prices(tickers_prices.keys())
            for ticker in tickers_bids_asks.keys():
                highest_bid, lowest_ask = tickers_bids_asks[ticker]
                if should_intervene(highest_bid, lowest_ask):
                    target_price = derive_price(highest_bid, lowest_ask)
                    price_data = (
                        highest_bid,
                        lowest_ask,
                        get_highest_bid(target_price),
                        get_lowest_ask(target_price)
                    )
                    intervene(market_maker_id, ticker, price_data)




def get_tickers_prices_from_file(filename):
    with open(filename) as file:
        reader = csv.reader(file)
        tickers_prices = {}
        for row in reader:
            tickers_prices[row[0]] = float(row[1])
        return tickers_prices


def parse_args(argv):
    parser = argparse.ArgumentParser(description='Market maker for trading trading application')
    parser.add_argument('-m', '--num_makers', type=int, help='Number of market makers')
    parser.add_argument('-t', '--tickers', type=str, help='Path to file that has tickers and prices')
    parser.add_argument('-l', '--logfile', type=str, help='Path to logfile')
    parser.add_argument('-d', '--debug', action='store_true', help='Enable DEBUG logging')

    return parser.parse_args(argv[1:])


def main(argv):
    parsed_args = parse_args(argv)

    logConfig = {
        'level': logging.DEBUG if parsed_args.debug else logging.INFO,
        'format': '%(asctime)s - %(name)s - %(levelname)s - %(filename)s:%(lineno)s - %(message)s'
    }

    if parsed_args.tickers is None:
        print 'No Tickers file found, exiting...'
        exit(1)

    if parsed_args.logfile:
        logConfig['filename'] = parsed_args.logfile
        logging.basicConfig(**logConfig)

    num_makers = parsed_args.num_makers if parsed_args.num_makers is not None else 1
    MARKET_MAKER_IDS = range(-1, -1 - num_makers, -1)

    tickers_dict = get_tickers_prices_from_file(parsed_args.tickers)

    try:
        start_market_makers(tickers_dict, MARKET_MAKER_IDS)
    except Exception as e:
        logging.exception(e)
        exit(1)


if __name__ == '__main__':
    import sys
    rc = main(sys.argv)
