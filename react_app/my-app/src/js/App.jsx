class App extends React.Component {
  render() {
    return (<Main />);
  }
}

class Main extends React.Component {

  constructor(props) {
    super(props);
    this.buy = this.buy.bind(this);
    this.sell = this.sell.bind(this);
    this.serverRequest = this.serverRequest.bind(this);
  }

  buy() {
    var dummy_data_buy = {"BuyerId" : 101, "SellerId" : 404, "Ticker" : "AAPL", "AmountTraded" : 42, "CashTraded" : 420};
    this.serverRequest(dummy_data_buy, "bid");
  }

  sell() {
    var dummy_data_sell = {"BuyerId" : 101, "SellerId" : 404, "Ticker" : "AAPL", "AmountTraded" : 42, "CashTraded" : 420};
    this.serverRequest(dummy_data_sell, "ask");
  }

  serverRequest(dummy_data_str, url_type) {
    var dummy_data = JSON.stringify(dummy_data_str);
    /* TODO: Change local host to the actual address of the server. */
    console.log("Sent POST request for request:" + url_type);
    $.post(
      "localhost:8080/api/" + url_type,
      dummy_data,
      res => {
        window.alert("Transaction completed at time:" + res);
        //console.log(res);
	console.log("Transaction done");
      }
    );
  }



  render() {
    return (
      <div className="main_stage">
        <button className="buy_button" onClick={this.buy}> Buy!!! </button>
        <button className="sell_button" onClick={this.sell}> Sell! </button>
      </div>
    )
  }
}

ReactDOM.render(<App />, document.getElementById('app'));
