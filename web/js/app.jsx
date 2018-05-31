class App extends React.Component {
  render() {
    return (<Main />);
  }
}

class Main extends React.Component {
  render() {
    return (
      <div className="main_stage">
        <button className="buy_button"> Buy </button>
        <button className="sell_button"> Sell </button>
      </div>
    )
  }
}