import React from 'react';
import ReactDOM from 'react-dom';
import $ from 'jquery';
import jQuery from 'jquery';
import './stylesheets/style.css';
import Select from 'react-select';
import 'react-select/dist/react-select.css';
import {LineChart} from 'react-easy-chart';

class App extends React.Component {

  componentWillMount() {
    window.MyVars = {
      //id: parseInt(prompt("What user ID?", "Enter user ID")),
      id: 1,
    }
  }

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
      selectValue: 'Apple Inc. (AAPL)',
      clearable: false,
      recentlyViewedList: [],
      options: [],
    }
    this.updateValue = this.updateValue.bind(this);
    this.jumpToRecent = this.jumpToRecent.bind(this);
  }

  componentDidMount() {

    var options_json = this.generateDummyOptions();
    
    var options = [];
    
    for (var i = 0; i < options_json.length; i++) {
      var name = "" + options_json[i].Label + " (" + options_json[i].Value + ")";
      console.log(name);
      options.push({value: name, label: name });
    }

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

  /*generateDummyOptions() {
    var result = [];
    var dummy_data_str = {"packet" : "hi"};
    var dummy_data = JSON.stringify(dummy_data_str);
    jQuery.ajaxSetup({async:false});
    $.get(
      "http://localhost:8080/api/get-company-list",
      //"http://cloud-vm-45-112.doc.ic.ac.uk:8080/api/get-company-list",
      res => {
       // console.log(res.results);
        result = res.results;
      }
    );
    return result; 
  }*/

  generateDummyOptions() {
    var result = [{Label: "Apple", Value: "AAPL"}, {Label: "Microsoft", Value: "MSFT"}];
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

//Function to get the initial data points, should call backend for this. 
function getInitialDataForGraph() {
  var initial_data = [];
  //initial_data.push({x: '1-Jan-15 10:00:00', y: 20});
  //initial_data.push({x: '1-Jan-15 10:00:30', y: 70});
  //initial_data.push({x: '1-Jan-15 10:01:00' , y: 40});

  const initialDate = '1-Jan-15 ';
 
  for(var i = 0; i < 10; i++) {

    var d = new Date(new Date().getTime() - ((10 - i) * 10 * 1000));
    var rnd = Math.floor((Math.random() * 80) + 20);
    
    var h = (d.getHours() < 10 ? '0' : '') + d.getHours();
    var m = (d.getMinutes() < 10 ? '0' : '') + d.getMinutes();
    var s = (d.getSeconds() < 10 ? '0' : '') + d.getSeconds();

    const dateToAdd = initialDate + h + ':' + m + ':' + s;

    initial_data.push({x: dateToAdd, y: rnd});
  }

  var wrapper = [];
  wrapper.push(initial_data);
  return wrapper;
}

// Function to get the next data point, should call backend for this. 
function getNextDataPointForGraph() {
  const initialDate = '1-Jan-15 ';
  // var rnd = Math.floor((Math.random() * 80) + 20);
  var d = new Date();
  var h = (d.getHours() < 10 ? '0' : '') + d.getHours();
  var m = (d.getMinutes() < 10 ? '0' : '') + d.getMinutes();
  var s = (d.getSeconds() < 10 ? '0' : '') + d.getSeconds();

  const dateToAdd = initialDate + h + ':' + m + ':' + s;
  
  // var nextPoint = {x: dateToAdd, y: rnd};
  return dateToAdd;
}

class Graph extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      graph_width: 0,
      graph_height: 0,
      data: [10], 
    };
    this.updateDataGraph = this.updateDataGraph.bind(this);
    this.myRef = React.createRef();
  }

  componentDidMount() {
    const boundingBox = this.myRef.current.getBoundingClientRect();
    const dataToPlot = getInitialDataForGraph();

    this.setState({
      graph_width: boundingBox.width,
      graph_height: boundingBox.height,
      data: dataToPlot,
    }, function() {
      this.updateDataGraph();
    }); 

    window.addEventListener("resize", this.updateDimensions.bind(this));

  }


  // Function to update the graph every 10 seconds with new data points. call backend for this. 
  updateDataGraph() {
  
    if (this.state.data.length === 0) {
      setTimeout(this.updateDataGraph, 1 * 1000);
    }
    const nextDataPointTime = getNextDataPointForGraph();
    var new_data_set = this.state.data[0].slice();
    new_data_set.shift();
    var rnd = Math.floor((Math.random() * 80) + 20);
    console.log(nextDataPointTime);
    new_data_set.push({x: nextDataPointTime, y: rnd});
    var wrapper = [];
    wrapper.push(new_data_set);
    this.setState({
      data: wrapper,
    });
    setTimeout(this.updateDataGraph, 10 * 1000);
  }


  updateDimensions() {
    const boundingBox = this.myRef.current.getBoundingClientRect();
    this.setState({
      graph_width: boundingBox.width,
      graph_height: boundingBox.height,
    }); 
  }

  render() {
    return (
      <div className="graph_display_cont">
        <div className="graph_display"> Showing graph for {this.props.current_company}:
        <div className="graph_cont" ref={this.myRef}>
        <LineChart
          datePattern={'%d-%b-%y %H:%M:%S'}
          xType={'time'}
          axes
          grid
          verticalGrid
          //interpolate={'cardinal'}
          lineColors={['cyan']}
          width={this.state.graph_width}
          height={this.state.graph_height}
          axisLabels={{x: 'Time', y: 'Price (USD)'}}
          data={this.state.data}
        />
        </div> 
        </div>
      </div>
    )
  }
}

function getFigures(comp) {
  
  var regExp = /\(([^)]+)\)/;
  var result = regExp.exec(comp);
  console.log(result)
  if (result === null) {
    return "";
  }
  var dummy_data_comp = {"UserId" : window.MyVars.id, "Ticker" : result[1]};
  var dummy_data = JSON.stringify(dummy_data_comp);
  var temp;
  jQuery.ajaxSetup({async:false});
  /*$.post(
    "http://localhost:8080/api/get-company-info",
    //"http://cloud-vm-45-112.doc.ic.ac.uk:8080/api/get-company-info",
    dummy_data,
    res => {
      temp = res.Amount;  
    }
  );*/
  console.log(temp)
  return temp;
}

class CompanyInfo extends React.Component {

  render() {

    var figures = getFigures(this.props.current_company);
    //var figures = [1, 2];
    console.log("updated");

    return (
      <div className="company_info_cont"> 
        Showing for {this.props.current_company}: 
        <br/> Price: {}$
        <br/> Currently own {figures} shares. 
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
  
  getTicker(str) {
    var regExp = /\(([^)]+)\)/;
    var result = regExp.exec(str);
    return result[1];
  }

  buy() {
    var thing_to_cut = this.props.current_company;
    var ticker = this.getTicker(thing_to_cut);
    console.log(ticker);
    var dummy_data_buy = {"userId" : window.MyVars.id, "equityTicker" : ticker, "amount" : 1, "orderType" : "marketBid"};
    this.serverRequest(dummy_data_buy, "bid");
  }

  sell() {
    var thing_to_cut = this.props.current_company;
    var ticker = this.getTicker(thing_to_cut);
    var dummy_data_sell = {"userId" : window.MyVars.id, "equityTicker" : ticker, "amount" : 1, "orderType" : "marketAsk"};
    this.serverRequest(dummy_data_sell, "ask");
  }

  serverRequest(dummy_data_str, url_type) {
    var dummy_data = JSON.stringify(dummy_data_str);
    /* TODO: Change local host to the actual address of the server. */
    console.log("Sent POST request for request:" + url_type);
    jQuery.ajaxSetup({async:false});
    /*$.post(
      //"http://cloud-vm-45-112.doc.ic.ac.uk:8080/api/" + url_type,
      "http://localhost:8080/api/" + url_type,
      dummy_data,
      res => {
        window.alert("Order submitted!");
      }
    );*/
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
          current_company={this.props.current_company}
        />
        <Button 
          button_type={"sell_button"}
          onClick={() => this.sell()}
          button_name={"Ask"}
          current_company={this.props.current_company}
        />
      </div>
    )
  }
}

/* Button class for rendering the buttons. */
class Button extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      renderChild: false,
    };
    this.handleChildUnmount = this.handleChildUnmount.bind(this);
    this.handleChildMount= this.handleChildMount.bind(this);
  }

  handleChildUnmount() {
    this.setState({renderChild: false});
  }

  handleChildMount() {
    this.setState({renderChild: true});
  }

  render() {
    return (
      <div className="button_and_action_wrapper">
        <button className={this.props.button_type} onClick={this.handleChildMount}> {this.props.button_name} </button> 
        {this.state.renderChild ? <ActionConfirmation unmountMe={this.handleChildUnmount} current_company={this.props.current_company} button_name={this.props.button_name}/> : null}
      </div> 
    )
  }
}

class ActionConfirmation extends React.Component {

  constructor(props) {
    super(props);

    this.state ={
      number_of_stock: 0,
      action_type: "market",
      sample_stock_value: 69,
      user_budget: 6942096,
      renderSubmitted: false,
    }
    this.handleCloseConfirmation = this.handleCloseConfirmation.bind(this);
  }

  dismiss() {
    this.props.unmountMe();
  }

  handleCloseConfirmation() {
    this.setState({renderSubmitted : false});
    this.dismiss();
  }

  inputChangeStock(e) {
    /* Disregards the decimal place, so always a whole number */ 
    var value = parseInt(e.target.value, 10);
    /* Always have a default of 0 even if input is empty. */
    if (Number.isNaN(value)) {
      value = 0;
    }

    this.setState({
      number_of_stock: value
    });
  }

  inputChangeAction(e) {
    const value = e.target.value;
    this.setState({
      action_type: value
    })
  }

  getTicker(str) {
    var regExp = /\(([^)]+)\)/;
    var result = regExp.exec(str);
    return result[1];
  }

  submitRequest() {
    var ticker = this.getTicker(this.props.current_company);
    var data_to_send ={"userId" : 1, "equityTicker" : ticker, "amount" : this.state.number_of_stock, "orderType" : this.state.action_type + this.props.button_name} ;
    var data = JSON.stringify(data_to_send);
    console.log(data);
    this.setState({renderSubmitted: true});
    /* 
    jQuery.ajaxSetup({async:false});
    var url_type = this.props.button_name.toLowerCase();
    $.post(
      //"http://cloud-vm-45-112.doc.ic.ac.uk:8080/api/" + url_type,
      "http://localhost:8080/api/" + url_type,
      dummy_data,
      res => {
        window.alert("Order submitted!");
      }
    );
    */
    // this.dismiss();
  }

  render() {

    var current_amount = this.state.sample_stock_value * this.state.number_of_stock;
    var amount_left = this.state.user_budget - current_amount;

    return (
      <div className="darken_bg">
        {this.state.renderSubmitted ? <SubmitConfirmation unmountMe={this.handleCloseConfirmation}/> : null}
        <div className="confirmation_window"> 
          <button className="close_button" onClick={() => this.dismiss()}> X </button> 
          <p className="company_viewing"> Viewing for {this.props.current_company} - {this.props.button_name}:</p>
          <p> Number of stock: <input type="number" onChange={e => this.inputChangeStock(e)}/> </p>
          <p> Type of action: 
            {/* <div> */}
              <input type="radio" onClick={e => this.inputChangeAction(e)} id="actionChoice1" name="action" value="market" defaultChecked/>
              <label htmlFor="actionChoice1"> Market </label>

              <input type="radio" onClick={e => this.inputChangeAction(e)} id="actionChoice2" name="action" value="limit"/>
              <label htmlFor="actionChoice1"> Limit </label>
            {/* </div> */}
          </p>
          <p> Total price: {current_amount}</p>
          <p style={{color: amount_left > 0 ? "black" : "red"}}> Total funds left: {amount_left}</p> 
          <div className="place_order">
            <button className="place_order_button" disabled={amount_left < 0 || this.state.number_of_stock <= 0} onClick={() => this.submitRequest()}> Place order </button> 
          </div>
        </div> 
      </div>
    )
  }
}

class SubmitConfirmation extends React.Component {

  dismiss() {
    this.props.unmountMe();
  }

  render() {
    return (
      <div className="darken_bg2">
        <div className="submit_window">
          <button className="ok_confirmation_button" onClick={() => this.dismiss()}> Submitted! <br /> Click to dismiss </button> 
        </div> 
      </div> 
    )
  }
}

ReactDOM.render(<App />, document.getElementById('root'));
