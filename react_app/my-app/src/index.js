import React from 'react';
import ReactDOM from 'react-dom';
import $ from 'jquery';
import './stylesheets/style.css';

class App extends React.Component {
  render() {
    //document.body.style.className = "c-container / t--light";
    return (
      <div id='Stage' className="c-container / t--light">
        <ThemeSelector 
          onClick={() => switch_theme()}
        />
        <UiInterface />
      </div>
    );
  }
}

function switch_theme() {
  console.log("Here");
  const body = document.getElementById('Stage');

  if (body.classList.contains('t--light')) {
    body.classList.remove('t--light');
    body.classList.add('t--dark');
  } else {
    body.classList.remove('t--dark');
    body.classList.add('t--light');
  }
}

function ThemeSelector(props) {
  return(
    <button className="ThemeSelector" onClick={props.onClick}> Switch Theme </button> 
  )
}

class UiInterface extends React.Component {

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
        <Button 
          button_type={"buy_button"}
          onClick={() => this.buy()}
          button_name={"Bid"}
        />
        <Button 
          button_type={"sell_button"}
          onClick={() => this.sell()}
          button_name={"Ask"}
        />
      </div>
    )
  }
}

/* Button class for rendering the buttons. */
function Button(props) {
  return (
    <button className={props.button_type} onClick={props.onClick}> {props.button_name} </button>
  )
}

ReactDOM.render(<App />, document.getElementById('root'));
