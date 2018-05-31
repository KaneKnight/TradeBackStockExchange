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
    var dummy_data_buy = {"userId" : 1};
    this.serverRequest(dummy_data_buy, "bid");
  }

  sell() {
    var dummy_data_sell = {"userId" : 1};
    this.serverRequest(dummy_data_sell, "ask");
  }

  serverRequest(dummy_data_str, url_type) {
    var dummy_data = JSON.stringify(dummy_data_str);
    $.post(
      "http://localhost:8080/api/" + url_type,
      dummy_data,
      res => {
        console.log("Success");
      }
    );
  }



  render() {
    return (
      <div className="main_stage">
        <button className="buy_button" onClick={this.buy}> Buy </button>
        <button className="sell_button" onClick={this.sell}> Sell </button>
      </div>
    )
  }
}

ReactDOM.render(<App />, document.getElementById('app'));