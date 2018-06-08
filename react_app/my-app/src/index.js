import React from 'react';
import ReactDOM from 'react-dom';
import $ from 'jquery';
import './stylesheets/style.css';
import Select from 'react-select';
import 'react-select/dist/react-select.css';

class App extends React.Component {
  render() {
    return (
      <Main />
    )
  }
}

class Main extends React.Component {

  constructor(props) {
    super(props);
    this.selectNewCompany = this.selectNewCompany.bind(this);
    this.state = {
      current_company: "Apple",
    };
  }

  selectNewCompany(new_company) {
    //console.log("Called with " + new_company);
    //console.log("Called");
    this.setState({current_company : new_company})
  }

  render() {
    //document.body.style.className = "c-container / t--light";
    return (
      //<div id='Stage' className="c-container / t--light">
      <div id='Stage' className="grid-container">
        <NavigationBar />
        <CompanyList onChange = {this.selectNewCompany}/>
        <GraphAndButtons onChange={this.selectNewCompany} current_company={this.state.current_company}/>
        <CompanyInfo current_company={this.state.current_company}/>
        <UserInfo />
      </div>
    );
  }
}

class CompanyList extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
			searchable: true,
      selectValue: 'test1',
      clearable: false,
      recentlyViewedList: [],
      options: [],
    }
    this.updateValue = this.updateValue.bind(this);
    this.jumpToRecent = this.jumpToRecent.bind(this);
  }

  componentDidMount() {

    var options = this.generateDummyOptions();
    this.setState({
      options: options,
    });

    this.props.onChange(this.state.selectValue);

  }

  updateValue(newValue) {
    this.setState({
			selectValue: newValue,
    });
    this.props.onChange(newValue);
    var new_list = this.state.recentlyViewedList.slice();

    /* Check if the key already exists, meaning that we have to move 
       it to the top of the list. */

    for (var i = 0; i < new_list.length; i++) {
      if(new_list[i].key === newValue) {
        new_list.splice(i, 1);
      }
    }

    /* Generate and insert the new element for the list. */
    var new_elem = <RecentElem key={newValue} recentId={newValue} jumpToRecent={this.jumpToRecent}/>;
    const max_number_of_elems_in_list = 5;
    new_list.unshift(new_elem);
    if (new_list.length > max_number_of_elems_in_list) {
      new_list.splice(max_number_of_elems_in_list, 1);
    }

    /* Update state list to the new generate one. */
    this.setState({
      recentlyViewedList: new_list,
    })
  }

  generateDummyOptions() {
    var result = [];
    for (var i = 0; i < 50; i++) {
      var name = 'test' + i;
      result.push({value: name, label: name});
    }
    return result;
  }

  jumpToRecent(elem) {
    var index_of_elem = -1;
    var new_list = this.state.recentlyViewedList.slice();

    /* Find index of the element we want to jump to*/
    for (var i = 0; i < new_list.length; i++) {
      if (new_list[i].key === elem) {
        index_of_elem = i;
        break;
      }
    }

    /* Remove it from original array. Ok to use original as updateValue will make a copy and then
       update it. */
    this.state.recentlyViewedList.splice(index_of_elem, 1);
    this.updateValue(elem);
  }

  render() {

    return (
      <div id='test1' className="company_list_cont">
        <div id='test2'>
        <Select
					id="state-select"
					ref={(ref) => { this.select = ref; }}
					onBlurResetsInput={false}
					onSelectResetsInput={false}
					autoFocus
					options={this.state.options}
					simpleValue
					name="selected-state"
					value={this.state.selectValue}
					onChange={this.updateValue}
          searchable={this.state.searchable}
          clearable={this.state.clearable}
				/>
        </div>
        <RecentlyViewed 
          recentlyViewedList={this.state.recentlyViewedList}
          jumpToRecent={this.jumpToRecent}
        />
      </div>
    )
  }
}

class RecentlyViewed extends React.Component {

  render() {

    return (
      <div id='recent_list' className='recently_viewed_cont'> 
        Recently viewed: 
        {this.props.recentlyViewedList}
      </div>
    )
  }
}

class RecentElem extends React.Component {

  doTheJump(elem) {
    this.props.jumpToRecent(elem);
  }

  render() {
    return (
      <div className='recent_elem'>
        <a id={this.props.recentId} href="javascript:;" onClick={() => this.doTheJump(this.props.recentId)}> {this.props.recentId} </a>
      </div>
    )
  }
}

class GraphAndButtons extends React.Component {
  render() {
    return (
      <div className="grid-container-graph"> 
        <Graph current_company={this.props.current_company}/>
        <UiInterface onChange={this.props.onChange} current_company={this.props.current_company}/>
      </div>
    )
  }
}

class Graph extends React.Component {

  constructor(props) {
    super(props);
    this.displayMessage = this.displayMessage.bind(this);
  }

  displayMessage() {
    window.alert(
      "Graph that shows the price history of a share for the current company."
    );
  }

  render() {
    return (
      <div className="graph_display_cont">
        <div className="graph_display"> Showing graph for {this.props.current_company}
        <button className="againPleaseStop" onClick={this.displayMessage}>?</button>
          <div className="photo">
          <img src="http://www.bbc.co.uk/staticarchive/d952b4abb2a9af3e8f001f7af8afaecfa7a4e4ae.gif" alt="dummy_graph_example" className="center"></img> 
          </div>
        </div>
      </div>
    )
  }
}

function getFigures(comp) {
  var dummy_data_comp = {"BuyerId" : 101, "company" : comp};
  var dummy_data = JSON.stringify(dummy_data_comp);
  var temp;
  $.get(
    "localhost:8080/api/get-company-info/",
    dummy_data,
    res => {
      temp = res;
    }
  );
  return temp;
}

class CompanyInfo extends React.Component {

  constructor(props) {
    super(props);
    this.displayMessage = this.displayMessage.bind(this);
  }

  displayMessage() {
    window.alert(
      "This represents the real-time price of the stock in the market. This is a very volatile figure, and can often change in minutes. This is due to a high volume of orders that get processed every second, which can rapidly change the price."
    );
  }

  render() {

    //var figures = getFigures(this.props.current_company);
    var figures = [1, 2];
    console.log("updated");

    return (
      <div className="company_info_cont"> 
        Showing for {this.props.current_company}: 
        <br/> Price: {figures[0]}$ <button className="yolowhocares" onClick={this.displayMessage}>?</button>
        <br/> Currently own {figures[1]} shares. 
      </div>
    )
  }
}

class UserInfo extends React.Component {
  render() {
    return (
      <div className="user_info_cont"> User Info Here </div>
    )
  }
}

class NavigationBar extends React.Component {
  render() {
    return (
      <div id='nav_bar' className="nav_bar_cont">
        <div id='grid_nav_bar' className="grid-container-nav-bar">
          <div className="app_name_cont"> App Name Here </div> 
          <div className="nav_gap_cont"> </div>
          <div className="theme_switch_cont"> Should be switch </div>
          <div className="login_btn_cont"> Login </div>
        </div>
      </div>
    )
  }
}

/* Function to switch between light and dark theme. */
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

/* Component to switch between light and dark theme. */
function ThemeSelector(props) {
  return(
    <div className="onoffswitch">
      <input type="checkbox" name="onoffswitch" className="onoffswitch-checkbox" id="myonoffswitch" defaultChecked onClick={props.onClick}/>
      <label className="onoffswitch-label" htmlFor="myonoffswitch">
        <span className="onoffswitch-inner"></span>
        <span className="onoffswitch-switch"></span>
      </label>
    </div>
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
    var dummy_data_buy = {"BuyerId" : 101, "SellerId" : 404, "Ticker" : this.props.current_company, "AmountTraded" : 42, "CashTraded" : 420};
    this.serverRequest(dummy_data_buy, "bid");
  }

  sell() {
    var dummy_data_sell = {"BuyerId" : 101, "SellerId" : 404, "Ticker" : this.props.current_company, "AmountTraded" : 42, "CashTraded" : 420};
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
        window.alert("Transaction completed!");
      }
    );
    /* Change to current company to update the view */
    this.props.onChange(this.props.current_company);
  }



  render() {
    return (
      <div className="ui_buttons_cont">
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
class Button extends React.Component {

  constructor(props) {
    super(props);
    this.displayMessage = this.displayMessage.bind(this);
  }

  displayMessage() {

    var msg = ""

    if (this.props.button_type === "buy_button") {
      msg = "There are two types of buy, a market buy and a limit buy:" + 
      "\n > Market Buy" + 
      "\n A market buy is an order to buy equity shares at the current available market price. The order will be completed at the current market price, assuming enough trading volume is available. Market buys should only be used for trades that need to happen quickly, with less priority given to price." +
      "\n > Limit Buy" + 
      "\n A limit buy is an order to buy equity shares at either the price specified by the investor or below. This should be most used type of buy. A limit buy ensures that the investor is only pays the certain maximum price for the shares, and is therefore only executed when this price is available."
    } else {
      msg = "There are two types of sell, a market sell and a limit sell" +
      "\n > Market Sell" + 
      "\n A market sell is an order to sell equity shares at the current available market price. This order will be filled at the current market price, provided there is enough trading volume to process a trade. Market sells should be used when the investor wants to sell quickly. This could be because either the price is rapidly dropping and the investor is trying to cut losses, or to reinvest that money in a better stock." + 
      "\n > Limit Sell" + 
      "\n An order to sell equity shares at a specified price, or better. Sell-limit orders can be used to specify a minimum price for which you are willing to sell equity shares, and will be executed only once that price (or a better one) is available. They can be useful if you have a target selling price in mind but are unable to frequently monitor your portfolio." +
      "This should be used by investors who aren’t in a hurry to sell, to try and get a better price. However if the market is incredibly volatile, often your order can be “leapfrogged,” such that it expires and the day could end with a lower price. So these should be used with caution."
    }
    window.alert(msg);
  }

  render() {
    return (
      <div className={this.props.button_type}>
      <button className={this.props.button_type + "_1"} onClick={this.props.onClick}> {this.props.button_name} </button>
      <button className="question_box" onClick={this.displayMessage}> ? </button>
      </div>
    )
  }
}

ReactDOM.render(<App />, document.getElementById('root'));
