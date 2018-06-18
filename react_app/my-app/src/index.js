import React from 'react';
import ReactDOM from 'react-dom';
import $ from 'jquery';
import jQuery from 'jquery';
import './stylesheets/style.css';
import Select from 'react-select';
import 'react-select/dist/react-select.css';
import {LineChart} from 'react-easy-chart';
import auth0 from 'auth0-js';
import 'react-interactive-tutorials/dist/react-interactive-tutorials.css';
import interactiveTutorials from 'react-interactive-tutorials';
import {
  paragraphs,
  registerTutorials,
  startTutorial
} from 'react-interactive-tutorials';

const AUTH0_CLIENT_ID = "pRPNKyXdY9dNx0yzmwhhrxi1PIDmHQ0v";
const AUTH0_DOMAIN = "ic22-webapps2018.eu.auth0.com";
const AUTH0_CALLBACK_URL = window.location.href;
const AUTH0_API_AUDIENCE = "webapps2018-tradingplt";

class App extends React.Component {

  parseHash() {
    this.auth0 = new auth0.WebAuth({
      domain: AUTH0_DOMAIN,
      clientID: AUTH0_CLIENT_ID
    });
    this.auth0.parseHash(window.location.hash, (err, authResult) => {
      if (err) {
        return console.log(err);
      }
      if (
        authResult !== null &&
        authResult.accessToken !== null &&
        authResult.idToken !== null
      ) {
        localStorage.setItem("access_token", authResult.accessToken);
        localStorage.setItem("id_token", authResult.idToken);
        localStorage.setItem(
          "profile",
          JSON.stringify(authResult.idTokenPayload)
        );
        window.location = window.location.href.substr(
          0,
          window.location.href.indexOf("#")
        );
      }
    });
  }

  setup() {
    $.ajaxSetup({
      beforeSend: (r) => {
        if (localStorage.getItem("access_token")) {
          r.setRequestHeader(
            "Authorization",
            "Bearer " + localStorage.getItem("access_token")
          );
        }
      }
    });
  }

  setState() {
    let idToken = localStorage.getItem("id_token");
    if (idToken) {
      this.loggedIn = true;
    } else {
      this.loggedIn = false;
    }
    // window.MyVars.id = 1;
  }

  componentWillMount() {
    this.setup();
    this.parseHash();
    this.setState();
  }

  render() {
    if (this.loggedIn) {
      return <Main />;
    }
    return <Home />;
  }

}

class Home extends React.Component {

  constructor(props) {
    super(props);
    this.authenticate = this.authenticate.bind(this);
  }
  authenticate() {
    this.WebAuth = new auth0.WebAuth({
      domain: AUTH0_DOMAIN,
      clientID: AUTH0_CLIENT_ID,
      scope: "openid profile",
      audience: AUTH0_API_AUDIENCE,
      responseType: "token id_token",
      redirectUri: AUTH0_CALLBACK_URL
    });
    this.WebAuth.authorize();
  }

  render() {
    return (
      <div className="landing_page"> 
      <div className="background_landing_page">
      </div>
      <div className="landing_page_text">
        Welcome to <span style={{fontWeight: "bold", color: "white"}}>TradeBack</span>, the sandbox platform to see if financial trading is for you! 
        <button className="sign_in_button" onClick={this.authenticate}> Sign In To Get Started! </button> 
        </div>
      </div> 
    )
  }
}

function sendExistingUserCheck() {
  console.log(JSON.parse(localStorage.getItem("profile")).sub);
  var string_data = {"userIdString" : JSON.parse(localStorage.getItem("profile")).sub};
  var data = JSON.stringify(string_data);
  $.post(
    "http://localhost:8080/api/check-user-exists",
    data,
    res => {
      console.log("Finished checking user");  
    }
  )
}

class Main extends React.Component {

  constructor(props) {
    super(props);
    this.selectNewCompany = this.selectNewCompany.bind(this);
    this.updateCurrentPrice = this.updateCurrentPrice.bind(this);
    this.setInitialPrice = this.setInitialPrice.bind(this);
    this.renderedNewGraph = this.renderedNewGraph.bind(this);
    this.logout = this.logout.bind(this);
    this.state = {
      current_company: "Apple Inc. (AAPL)",
      current_price: -1,
      is_price_up: null,
      temp_price_history: [],
      need_to_update_graph: true,
    };
    sendExistingUserCheck();
    // $.post({
    //   "http://localhost:8080/api/check-existing-user",
    //   JSON.stringify(JSON.parse(localStorage.getItem("")))
    // });
  }

  logout() {
    localStorage.removeItem("id_token");
    localStorage.removeItem("access_token");
    localStorage.removeItem("profile");
    window.location.reload();
  }

  updateCurrentPrice(new_price) {
    
    //To make the inital display text not show down or up. 
    if (this.state.current_price === -1) {
      //Push to history so can later compare the values for initial pricing. 
      this.state.temp_price_history.push(new_price);
      this.setState({
        current_price: new_price,
      })
      return;
    }

    var is_new_price_up;
     
    //Check if there is a backed up queue. massive hack, don't look. 
    if (this.state.temp_price_history.length !== 0) {
      
      const len = this.state.temp_price_history.length;

      const to_compare = this.state.temp_price_history[len - 1];

      is_new_price_up = to_compare < new_price;

      //empty the temp queue.
      this.state.temp_price_history = [];
    } else {
      is_new_price_up = this.state.current_price < new_price;
    }

    // const is_new_price_up = this.state.current_price < new_price;
    // console.log("Comparing " + this.state.current_price + ' and ' + new_price);
   
    this.setState({
      current_price: new_price,
      is_price_up: is_new_price_up,
    });
  }

  setInitialPrice(initial_price) {
    this.setState({
      current_price: initial_price,
    });
  }

  selectNewCompany(new_company) {
    //console.log("Called with " + new_company);
    //console.log("Called");
    this.setState({
      current_company : new_company,
      need_to_update_graph : true, 
    });
  }

  renderedNewGraph() {
    this.setState({
      need_to_update_graph : false,
    });
  }

  render() {

    // console.log(localStorage.getItem("id_token"));
    // console.log(JSON.parse(localStorage.getItem("profile")).sub);

    //document.body.style.className = "c-container / t--light";
    return (
      //<div id='Stage' className="c-container / t--light">
      <div id='Stage' className="grid-container">
        <NavigationBar logout={this.logout}/>
        <CompanyList onChange = {this.selectNewCompany}/>
        <GraphAndButtons onChange={this.selectNewCompany} current_company={this.state.current_company} onPriceUpdate={this.updateCurrentPrice} current_price={this.state.current_price} setInitialPrice={this.setInitialPrice} need_to_update_graph={this.state.need_to_update_graph} renderedNewGraph={this.renderedNewGraph}/>
        <CompanyInfo current_company={this.state.current_company} current_price={this.state.current_price} is_price_up={this.state.is_price_up}/>
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
    
    //Change default text for drop down menu. 
    $(".Select-placeholder").html("Select a company...");

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
    var dummy_data_str = {"packet" : "hi"};
    var dummy_data = JSON.stringify(dummy_data_str);
    jQuery.ajaxSetup({async:false});
    $.get(
      "http://localhost:8080/api/get-company-list",
      //"http://cloud-vm-45-112.doc.ic.ac.uk:8080/api/get-company-list",
      res => {
        result = res.results;
      }
    );
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
        {/* <div className="graph_display_text"></div> */}
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
        <p className="indent_recent_title"> Recently viewed: </p>
        <ul> 
          {this.props.recentlyViewedList}
        </ul>
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
        <li> <a id={this.props.recentId} href="javascript:;" onClick={() => this.doTheJump(this.props.recentId)}> {this.props.recentId} </a> </li> 
      </div>
    )
  }
}

class GraphAndButtons extends React.Component {
  render() {
    return (
      <div className="grid-container-graph"> 
        <Graph current_company={this.props.current_company} current_price={this.props.current_price} onPriceUpdate={this.props.onPriceUpdate} setInitialPrice={this.props.setInitialPrice} need_to_update_graph={this.props.need_to_update_graph} renderedNewGraph={this.props.renderedNewGraph}/>
        <UiInterface onChange={this.props.onChange} current_company={this.props.current_company} current_price={this.props.current_price}/>
      </div>
    )
  }
}

//Function to get the initial data points, should call backend for this. 
function getInitialDataForGraph(comp) {
  var initial_data = [];
  //initial_data.push({x: '1-Jan-15 10:00:00', y: 20});
  //initial_data.push({x: '1-Jan-15 10:00:30', y: 70});
  //initial_data.push({x: '1-Jan-15 10:01:00' , y: 40});

  const initialDate = '18-Jun-18 ';
  /*
  for(var i = 0; i < dataPoints; i++) {

    var d = new Date(new Date().getTime() - ((dataPoints - i) * 10 * 1000));
    var rnd = Math.floor((Math.random() * 80) + 20);
    
    var h = (d.getHours() < 10 ? '0' : '') + d.getHours();
    var m = (d.getMinutes() < 10 ? '0' : '') + d.getMinutes();
    var s = (d.getSeconds() < 10 ? '0' : '') + d.getSeconds();

    const dateToAdd = initialDate + h + ':' + m + ':' + s;

    initial_data.push({x: dateToAdd, y: rnd});
  }*/

  var regExp = /\(([^)]+)\)/;
  var result = regExp.exec(comp);

  const data_to_send = {"ticker" : result[1], "dataNums" : 1};
  const data = JSON.stringify(data_to_send);
  var temp; 
  
  jQuery.ajaxSetup({async:false});
  $.post(
    "http://localhost:8080/api/get-datapoints",
    //"http://cloud-vm-45-112.doc.ic.ac.uk:8080/api/get-datapoints",
    data,
    res => {
      console.log(res.data[0].Price);
      temp = res.data[0].Price;
    }
  );

  var d = new Date();
  var h = (d.getHours() < 10 ? '0' : '') + d.getHours();
  var m = (d.getMinutes() < 10 ? '0' : '') + d.getMinutes();
  var s = (d.getSeconds() < 10 ? '0' : '') + d.getSeconds();

  const dateToAdd = initialDate + h + ':' + m + ':' + s;
  if (temp === undefined) {
    temp = 420;
  }

  initial_data.push({x: dateToAdd, y: temp});
  
  var wrapper = [];
  wrapper.push(initial_data);
  return wrapper;
}

// Function to get the next data point, should call backend for this. 
function getNextDataPointForGraph(comp) {
  const initialDate = '18-Jun-18 ';
  // var rnd = Math.floor((Math.random() * 80) + 20);
  var d = new Date();
  var h = (d.getHours() < 10 ? '0' : '') + d.getHours();
  var m = (d.getMinutes() < 10 ? '0' : '') + d.getMinutes();
  var s = (d.getSeconds() < 10 ? '0' : '') + d.getSeconds();

  const dateToAdd = initialDate + h + ':' + m + ':' + s;
  
  var regExp = /\(([^)]+)\)/;
  var result = regExp.exec(comp);

  const data_to_send = {"ticker" : result[1], "dataNums" : 1};
  const data = JSON.stringify(data_to_send);
  var temp; 
  
  jQuery.ajaxSetup({async:false});
  $.post(
    "http://localhost:8080/api/get-datapoints",
    //"http://cloud-vm-45-112.doc.ic.ac.uk:8080/api/get-datapoints",
    data,
    res => {
      temp = res.data[0].Price;  
    }
  );

  if (temp === undefined) {
    temp = 420;
  }

  var nextPoint = {x: dateToAdd, y: temp};
  return nextPoint;
}

var timeOutVar;

class Graph extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      graph_width: 0,
      graph_height: 0,
      dataPoints: 10,
      data: [], 
    };
    this.updateDataGraph = this.updateDataGraph.bind(this);
    this.myRef = React.createRef();
    this.updateToDifferentView = this.updateToDifferentView.bind(this); 
  }

  componentDidMount() {
    const boundingBox = this.myRef.current.getBoundingClientRect();
    const dataToPlot = getInitialDataForGraph(this.props.current_company);

    this.setState({
      graph_width: boundingBox.width,
      graph_height: boundingBox.height,
      data: dataToPlot,
    }, function() {
      const recent_price = dataToPlot[0][0].y;
      const newest_value = dataToPlot[0][0].y;
      this.props.setInitialPrice(recent_price);
      this.props.onPriceUpdate(newest_value);
      this.updateDataGraph();
    }); 

    window.addEventListener("resize", this.updateDimensions.bind(this));

  }

  //Todo: fix up/down not being correctly updated. Might have to flush the old value. 
  componentDidUpdate() {
    if (this.props.need_to_update_graph) {
      const newDataToPlot = getInitialDataForGraph(this.props.current_company);
      this.setState({
        data: newDataToPlot,
      }, function() {
        var recentIndex = this.state.data.length > 1 ? this.state.data.length - 2 : 0 ;
        const recent_price = newDataToPlot[0][recentIndex].y;
        const newest_value = newDataToPlot[0][this.state.data.length - 1].y;
        this.props.setInitialPrice(recent_price);
        this.props.onPriceUpdate(newest_value);
        //this.updateDataGraph();
      });
      this.props.renderedNewGraph();
    }
  }


  // Function to update the graph every 10 seconds with new data points. call backend for this. 
  //Set default for function to false; 
  updateDataGraph(stopUpdate = false) {

    if(stopUpdate) {
      clearTimeout(timeOutVar);
    } else {
  
    if (this.state.data.length === 0) {
      setTimeout(this.updateDataGraph, 1 * 1000);
    }
    const nextDataPointTime = getNextDataPointForGraph(this.props.current_company);
    var new_data_set = this.state.data[0].slice();
    if (new_data_set.length === 10) {
      new_data_set.shift();
    }
    // var rnd = Math.floor((Math.random() * 80) + 20);
    console.log(nextDataPointTime);
    new_data_set.push(nextDataPointTime);
    var wrapper = [];
    wrapper.push(new_data_set);
    this.setState({
      data: wrapper,
    });
    //Update the price to be the newly generated price for the display on the side. 
    this.props.onPriceUpdate(nextDataPointTime.y);
    //10 is the timeout time in amounts of seconds.
    //TODO: save the settimeout as a var and then clearTimeout(var) to stop it.  
    timeOutVar = setTimeout(this.updateDataGraph, 5 * 1000);
  }
  }

  updateToDifferentView() {
    var newDataPoints;
    if (this.state.dataPoints === 10) {
      newDataPoints = 30;
    } else {
      newDataPoints = 10; 
    }
    const newDataToPlot = getInitialDataForGraph(newDataPoints);
    this.setState({
      dataPoints: newDataPoints,
      data: newDataToPlot,
    }, function() {
      const recent_price = newDataToPlot[0][newDataPoints - 2].y;
      const newest_value = newDataToPlot[0][newDataPoints - 1].y;
      this.props.setInitialPrice(recent_price);
      this.props.onPriceUpdate(newest_value);
      this.updateDataGraph(newDataPoints === 30);
    });
    this.props.renderedNewGraph();
  }


  updateDimensions() {
    const boundingBox = this.myRef.current.getBoundingClientRect();
    this.setState({
      graph_width: boundingBox.width,
      graph_height: boundingBox.height,
    }); 
  }

  render() {

    // if (this.props.need_to_update_graph) {
    //   console.log("Should update to a new graph");
    //   this.props.renderedNewGraph();
    // }

    return (
      <div className="graph_display_cont">
        <div className="graph_display"> Showing graph for {this.props.current_company}: <button className="changeToMonth_button" onClick={this.updateToDifferentView}> {this.state.dataPoints === 10 ? "View Month" : "View Day"} </button> 
        <div className="graph_cont" ref={this.myRef}>
        <LineChart
          datePattern={'%d-%b-%y %H:%M:%S'}
          xType={'time'}
          axisLabels={{x: 'Time', y: 'Price (USD)'}}
          yDomainRange={[this.props.current_price - 2, this.props.current_price + 2]}
          axes
          grid
          verticalGrid
          // interpolate={'basis'}
          lineColors={['cyan']}
          //'#e9ecef',
          style={{
            // 'font-size' : '60px',
            'background-color': '#272B30',
            // '.tick line': {
            //   stroke: 'red',
            // },
            '.axis' : {
              stroke: 'white',
              // strokeWidth: 1,
              fill: 'white'
            },
            '.line' : {
              strokeWidth: 4,
            }
          }}
          width={this.state.graph_width}
          height={this.state.graph_height}
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
  var dummy_data_comp = {"userIdString" : JSON.parse(localStorage.getItem("profile")).sub, "Ticker" : result[1]};
  var dummy_data = JSON.stringify(dummy_data_comp);
  var temp;
  jQuery.ajaxSetup({async:false});
  $.post(
    "http://localhost:8080/api/get-company-info",
    //"http://cloud-vm-45-112.doc.ic.ac.uk:8080/api/get-company-info",
    dummy_data,
    res => {
      temp = res.Amount;  
    }
  );
  console.log('Figures returned was ' + temp);
  return temp;
}

class CompanyInfo extends React.Component {

  // constructor(props) {
  //   super(props);
  //   this.state = {
  //     number_of_renders: 0,
  //   }
  // }

  // componentDidUpdate() {
  //   const is_new_price_up = this.state.old_price < this.props.current_price;
  //   //console.log(this.state.old_price);
  //   //console.log(this.props.new_price);
  //   console.log(is_new_price_up);
  //   if (this.state.old_price !== this.props.current_price) {
  //     this.setState({
  //      old_price: this.props.current_price,
  //     })
  //   }
  // }

  render() {

    var figures = getFigures(this.props.current_company);
    //var figures = [1, 2];
    // console.log(this.state.number_of_renders);

    //TODO: Implement figures === 0 -> 'no shares' 
    console.log(figures);

    return (
      <div className="company_info_cont"> 
        <div className="company_info_text">
          Showing for <u> {this.props.current_company} </u>: 
          <p> Price: {this.props.current_price} $ {this.props.is_price_up === null ? null : (this.props.is_price_up ? '(up)' : '(down)')} </p>
          <p> Currently own {figures === undefined ? "no" : figures} shares. </p>
        </div>
      </div>
    )
  }
}

//TODO: start with this
class UserInfo extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      renderFullProfile: false,
    }

    this.handleRenderFullProfile = this.handleRenderFullProfile.bind(this);
    this.handleUnrenderFullProfile = this.handleUnrenderFullProfile.bind(this);

  }

  handleRenderFullProfile() {
    this.setState({
      renderFullProfile: true,
    });
  }

  handleUnrenderFullProfile() {
    this.setState({
      renderFullProfile: false, 
    });
  }


  render() {
    return (
      <div className="user_info_cont">
        <div className="user_info_area">
          <div className="info_overview">
            <span style={{fontSize: "20px", fontWeight: "bold"}}>User Portfolio preview: </span> 
            <p> Equities owned: </p>
            <p> Portfolio value: </p>  
          </div> 
          <div className="profile_button_cont">
            <button className="view_full_profile_button" onClick={this.handleRenderFullProfile}> Click To View Full Profile </button>  
            {this.state.renderFullProfile ? <FullUserProfile unmountMe={this.handleUnrenderFullProfile}/> : null}
          </div>
        </div>
      </div> 
    )
  }
}

function getPositionsForUser() { 

  //UNCOMMENT ME TO RUN IN NORMAL MODE 
  
  var positions = [];

  var data_to_send = {"userIdString": JSON.parse(localStorage.getItem("profile")).sub}
  var data = JSON.stringify(data_to_send);
  jQuery.ajaxSetup({async:false});
  $.post(
    "http://localhost:8080/api/get-all-user-positions",
    data,
    res => {
      positions = res;
    }
  )
  
  /*
  var positions = {
    "positions" : [
      {
        "ticker" : "AAPL",
        "numberOfSharesOwned" : 10,
        "valueOfPosition" : 1467,
        "percentageGain" : 1.6,
        "name" : "Apple"
      },
      {
        "ticker" : "MSFT",
        "numberOfSharesOwned" : 20,
        "valueOfPosition" : 14237,
        "percentageGain" : 10.6,
        "name" : "Microsoft"
      },
      {
        "ticker" : "BLZD",
        "numberOfSharesOwned" : 5,
        "valueOfPosition" : 762,
        "percentageGain" : -8,
        "name" : "Blizzard"
      },
    ]
  }*/

  return positions.positions;
}

function getTransactionHistoryForUser() {
  //Uncomment me to run with api call 
  
  var transactionHistory = [];
  var data_to_send = {"userIdString": JSON.parse(localStorage.getItem("profile")).sub}
  var data = JSON.stringify(data_to_send);
  jQuery.ajaxSetup({async: false});
  $.post(
    "http://localhost:8080/api/get-transaction-history",
    data,
    res => {
      transactionHistory[0] = res.BuyTransactions;
      transactionHistory[1] = res.SellTransactions;
    }
  )
  console.log(transactionHistory);
  
  
  /*const transactionHistory = {
    "BuyTransactions" : [
      {
        "ticker": "AAPL",
        "amountTraded": 1,
        "cashSpent" : 100,
        "price" : 100,
        "time" : new Date(new Date() - (Math.random() * 10000000) + 1),
      },
      {
        "ticker": "MSFT",
        "amountTraded": 4,
        "cashSpent" : 600,
        "price" : 150,
        "time" : new Date(new Date() - (Math.random() * 10000000) + 1),
      },
      {
        "ticker": "BLZD",
        "amountTraded": 3,
        "cashSpent" : 60,
        "price" : 20,
        "time" : new Date(new Date() - (Math.random() * 10000000) + 1),
      },
    ],
    "SellTransactions" : [
      {
        "ticker": "AAPL",
        "amountTraded": 1,
        "cashSpent" : 100,
        "price" : 100,
        "time" : new Date(new Date() - (Math.random() * 10000000) + 1),
      },
      {
        "ticker": "MSFT",
        "amountTraded": 4,
        "cashSpent" : 600,
        "price" : 150,
        "time" : new Date(new Date() - (Math.random() * 10000000) + 1),
      },
      {
        "ticker": "BLZD",
        "amountTraded": 3,
        "cashSpent" : 60,
        "price" : 20,
        "time" : new Date(new Date() - (Math.random() * 10000000) + 1),
      },
    ],
  }*/

  return transactionHistory;
}

class FullUserProfile extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      positions: getPositionsForUser(),
      transactionHistory: getTransactionHistoryForUser(),
    }
  }

  // componentDidMount() {
  //   var positions_of_user = getPositionsForUser();
  //   console.log(positions_of_user);
  //   this.setState({
  //     positions: positions_of_user,
  //   });
  // }

  render() {

    const to_stringify = {Name: "Louis Carteron", Current_amount: 1337, Starting_amount: 1000};

    //const user_profile = JSON.stringify(to_stringify);

    // console.log(user_profile);
    // console.log(to_stringify.Name);

    const price_difference = to_stringify.Current_amount - to_stringify.Starting_amount; 

    return (
      <div className="fake_new_page_bg">
        <div className="full_user_profile_wrapper"> 
          <button className="close_user_profile_button" onClick={this.props.unmountMe}>X</button> 
          <p style={{textAlign: "center"}}> Your User Portfolio: </p>
          <div className="user_info_profile_wrapper">
          <p style={{textAlign: "center"}}> Current Amount: {to_stringify.Current_amount} USD (<span style={{color: price_difference >= 0 ? "#53be53" : "#ee5f5b"}}>{price_difference >= 0 ? "Gained" : "Lost"} </span> {Math.abs(price_difference)} USD)</p> 
          </div> 
          <p style={{textAlign: "center", textDecoration: "underline"}}> Positions Owned: </p> 
          <div className="positions_held_wrapper">
            <Positions listOfPositions={this.state.positions}/>
          </div> 
          <br/> 
          <p style={{textAlign: "center", textDecoration: "underline"}}> Transaction History: </p>
          <div className="exchange_history_wrapper">
            <TransactionHistory transactionHistory={this.state.transactionHistory} />
          </div> 
        </div>
      </div> 
    )
  }
}

class Positions extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      positionsToShow: [],
    };
  };

  componentDidMount() {

    console.log(this.props.listOfPositions);

    var newPosToShow = [];
    for (var i = 0; this.props.listOfPositions !== null && i < this.props.listOfPositions.length; i++) {
      // var text_to_show = <p> Ticker: {this.props.listOfPositions.Positions[i].Ticker} Amount: {this.props.listOfPositions.Positions[i].Amount} Cash Spent: {this.props.listOfPositions.Positions[i].CashSpentOnPosition}$</p>
      var text_to_show = <div className="position_list_elem"> Company : {this.props.listOfPositions[i].name} ({this.props.listOfPositions[i].ticker}), <br/> Number of shares owned: {this.props.listOfPositions[i].numberOfSharesOwned}, <br/> Value of Position: {this.props.listOfPositions[i].valueOfPosition} USD, <br/> Percentage Gain: {this.props.listOfPositions[i].percentageGain}% </div>
      newPosToShow.push(text_to_show);
    }
    this.setState({
      positionsToShow: newPosToShow,
    })
  }

  render() {
    return (
      <div className="list_of_positions">
        {this.state.positionsToShow}
      </div>
    )
  }
}

class TransactionHistory extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      transactionHistoryToShow: [],
    };
  };

  componentDidMount() {
    var historyToShow = [];

    const buy_text = <p style={{color: "green"}}> Your Buy Transaction History: </p>;
    historyToShow.push(buy_text);
    //Looping through the Buy Transactions
    for (var i = 0; this.props.transactionHistory[0] !== null && i < this.props.transactionHistory[0].length; i++) {
      var text_to_show_for_buy = <div className="transaction_list_elem"> Ticker: {this.props.transactionHistory[0][i].ticker} <br/>
        Amount Traded: {this.props.transactionHistory[0][i].amountTraded}<br/>
        Cash Spent: {this.props.transactionHistory[0][i].cashSpent} USD<br/>
        Price Bought At: {this.props.transactionHistory[0][i].price} USD<br/>
        Transaction Time: {this.props.transactionHistory[0][i].time}</div>
      historyToShow.push(text_to_show_for_buy);
    }

    const sell_text = <p style={{color: "red"}}> Your Sell Transaction History: </p>;
    historyToShow.push(sell_text);
    //Looping through the Sell Transactions
    for (var i = 0; this.props.transactionHistory[1] !== null && i < this.props.transactionHistory[1].length; i++) {
      var text_to_show_for_buy = <div className="transaction_list_elem"> Ticker: {this.props.transactionHistory[1][i].ticker} <br/>
        Amount Traded: {this.props.transactionHistory[1][i].amountTraded} <br/>
        Cash Spent: {this.props.transactionHistory[1][i].cashSpent} USD<br/>
        Price Bought At: {this.props.transactionHistory[1][i].price} USD<br/>
        Transaction Time: {this.props.transactionHistory[1][i].time}</div>
      historyToShow.push(text_to_show_for_buy);
    }

    this.setState({
      transactionHistoryToShow: historyToShow,
    })
  }

  render() {
    return (
      <div className="transaction_history">
        {this.state.transactionHistoryToShow}
      </div>
    )
  }
}

const TUTORIALS = {
  'demo': {
    key: 'demo',
    title: 'TradeBack Tutorial',
    steps: [
      {
        key: 'anywhere',
        announce: paragraphs`
          This tutorial will explain how to use TradeBack so you can dive straight in 
          and experiment with trading, to figure out if it's for you!
        `,
        announceDismiss: "Okay, let's begin",
        activeWhen: [],
      },
      {
        key: 'beginning',
        highlight: '#test2',
        announce: paragraphs`
          This is where you can choose which company's stock you want to trade. 
        `,
        // annotateIn: 'div#test2 > div',
        // annotateSkip: 'Skip',
        announceDismiss: "Next",
        activeWhen: [
          {
            compare: 'checkpointComplete',
            checkpoint: 'demo_anywhere',
          },
        ],
      },
      {
        key: 'graph',
        highlight: '.graph_display',
        annotate: paragraphs`
          This graph shows the current price and the history for the company's
          stock that you selected.  
        `,
        annotateIn: 'div.graph_display_cont > div',
        // annotateSkip: 'Skip',
        annotateSkip: "Next",
        activeWhen: [
          {
            compare: 'checkpointComplete',
            checkpoint: 'demo_beginning',
          },
        ],
      },
      {
        key: 'view_month',
        highlight: '.changeToMonth_button',
        announce: paragraphs`
          Click this button to view the price history for the past month.  
        `,
        // annotateIn: 'div#test2 > div',
        // annotateSkip: 'Skip',
        announceDismiss: "Next",
        activeWhen: [
          {
            compare: 'checkpointComplete',
            checkpoint: 'demo_graph',
          },
        ],
      },
      {
        key: 'buy_button',
        highlight: '.buy_button',
        announce: paragraphs`
          This is where you can place a bid order to buy some stock for the company
          that you selected. 
        `,
        // annotateIn: 'div#test2 > div',
        // annotateSkip: 'Skip',
        announceDismiss: "Next",
        activeWhen: [
          {
            compare: 'checkpointComplete',
            checkpoint: 'demo_view_month',
          },
        ],
      },
      {
        key: 'sell_button',
        highlight: '.sell_button',
        announce: paragraphs`
          This is where you can place an ask order to sell some stock that you
          own for the company that you selected. 
        `,
        // annotateIn: 'div#test2 > div',
        // annotateSkip: 'Skip',
        announceDismiss: "Next",
        activeWhen: [
          {
            compare: 'checkpointComplete',
            checkpoint: 'demo_buy_button',
          },
        ],
      },
      {
        key: 'company_info',
        highlight: '.company_info_cont',
        announce: paragraphs`
          Here you can view usefull information about the company you selected, namely
          what the current market price for a share in that company is and how much
          stock you already own for that company.
        `,
        announceDismiss: "Next",
        activeWhen: [
          {
            compare: 'checkpointComplete',
            checkpoint: 'demo_sell_button',
          }
        ]
      },
      {
        key: 'user_info',
        highlight: '.user_info_cont',
        announce: paragraphs`
          Here you can view a short preview of your user portfolio.
        `,
        announceDismiss: "Next",
        activeWhen: [
          {
            compare: 'checkpointComplete',
            checkpoint: 'demo_company_info',
          }
        ]
      },
      {
        key: 'user_info_button',
        highlight: '.view_full_profile_button',
        announce: paragraphs`
          You can click this button to view your full portfolio, that will show 
          your current portfolio value, positions owned and your transaction history.
        `,
        announceDismiss: "Next",
        activeWhen: [
          {
            compare: 'checkpointComplete',
            checkpoint: 'demo_user_info',
          }
        ]
      },
      {
        key: 'complete',
        highlight: '.buy_button_popup_button, .sell_button_popup_button',
        announce: paragraphs`
          You will see these buttons scattered throughout the application. Don't be scared!
          If you feel like you don't understand the terminology, simple click this button 
          and you will be preseneted with an explanation!
        `,
        announceDismiss: "Next",
        activeWhen: [
          {
            compare: 'checkpointComplete',
            checkpoint: 'demo_user_info_button',
          }
        ]
      },
    ],
    complete: {
      on: 'checkpointReached',
      checkpoint: 'complete',
      title: 'TradeBack Tutorial Complete!',
      message: paragraphs`
        You have finished the tutorial!

        Now the market is yours. We have given you an initial amount of 10 Million USD.
        You know the basics, now it's time to start experimenting and learning. 
        Trade at will, and have fun! 
      `,
    },
  },
};

registerTutorials(TUTORIALS);

class NavigationBar extends React.Component {

  constructor(props) {
    super(props);
    this.launchTutorial = this.launchTutorial.bind(this);
  }

  launchTutorial() {
    console.log("called");
    startTutorial('demo');
    console.log("Finished");
  }

  
  render() {
    return (
      <div id='nav_bar' className="nav_bar_cont">
        <div id='grid_nav_bar' className="grid-container-nav-bar">
          <div className="app_name_cont"> TradeBack </div> 
          <div className="nav_gap_cont"> </div>
          <div className="theme_switch_cont"> <button className="view_tutorial_button" onClick={this.launchTutorial}> View Tutorial </button> </div>
          <div className="login_btn_cont"> <button className="logout_button" onClick={this.props.logout}> Logout </button> </div>
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
    var dummy_data_buy = {"userId" : 1, "equityTicker" : ticker, "amount" : 1, "orderType" : "marketBid"};
    this.serverRequest(dummy_data_buy, "bid");
  }

  sell() {
    var thing_to_cut = this.props.current_company;
    var ticker = this.getTicker(thing_to_cut);
    var dummy_data_sell = {"userId" : 1, "equityTicker" : ticker, "amount" : 1, "orderType" : "marketAsk"};
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
          current_price={this.props.current_price}
        />
        <Button 
          button_type={"sell_button"}
          onClick={() => this.sell()}
          button_name={"Ask"}
          current_company={this.props.current_company}
          current_price={this.props.current_price}
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
      renderInfoBubble: false,
      textForBuy: "A bid is a market offer to buy an amount of stock you want to own.",
      textForSell: "An ask is a market offer to sell an amount of stock you own.",
    };
    this.handleChildUnmount = this.handleChildUnmount.bind(this);
    this.handleChildMount = this.handleChildMount.bind(this);
    this.handlePopupMount = this.handlePopupMount.bind(this);
    this.handlePopupUnmount =  this.handlePopupUnmount.bind(this);
  }

  handleChildUnmount() {
    this.setState({renderChild: false});
  }

  handleChildMount() {
    this.setState({renderChild: true});
  }

  handlePopupUnmount() {
    this.setState({renderInfoBubble: false});
  }

  handlePopupMount() {
    this.setState({renderInfoBubble: !this.state.renderInfoBubble});
  }

  render() {

    return (
      // <div className="tempIdea">
      <div className={this.props.button_type + "_wrapper"}>
        
        <button className={this.props.button_type} onClick={this.handleChildMount}> {this.props.button_name} </button> 
        {this.state.renderChild ? <ActionConfirmation unmountMe={this.handleChildUnmount} current_company={this.props.current_company} button_name={this.props.button_name} current_price={this.props.current_price}/> : null}
        
        <div className="temp_idea">
          <button className={this.props.button_type + "_popup_button"} title="What is this?" onClick={this.handlePopupMount}> ? </button>
          <div className="temp_idea2">
            <div className="temp_idea3">
              {this.state.renderInfoBubble ? <InfoBubble text={this.props.button_type === "buy_button" ? this.state.textForBuy : this.state.textForSell}/> : null}
            </div>
          </div>
        </div> 
        {/* <InfoBubble /> */}
      {/* </div>  */}
      {/* <button> ? </button>  */}
      </div>
    )
  }
}

class InfoBubble extends React.Component {
  render() {
    return (
      <div className="info_bubble_wrapper">
        <div className="speech-bubble"> 
          <p>
          {this.props.text} 
          </p>
        </div> 
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
      help_type: "market",
      sample_stock_value: 69,
      user_budget: 6942096,
      renderSubmitted: false,
      limit_price: 0,
      showHelp: false,
      helpBidLimit: "A limit buy is an order to buy equity shares at either the price specified by the investor or below. This should be the most used type of buy. A limit buy ensures that the investor only pays a certain maximum price for the shares, and is therefore only executed when this price is available. ",
      helpBidMarket: "A market buy is an order to buy equity shares at the current available market price. The order will be completed at the current market price, assuming enough trading volume is available. Market buys should ONLY be used for trades that need to happen quickly, with less priority given to price.",
      helpAskLimit: "An order to sell equity shares at a specified price, or better. Sell-limit orders can be used to specify a minimum price for which you are willing to sell equity shares, and will be executed only once that price (or a better one) is available. They can be useful if you have a target selling price in mind but are unable to frequently monitor your portfolio. This should be used by investors who aren’t in a hurry to sell, to try and get a better price. However if the market is incredibly volatile, often your order can be “leapfrogged,” such that it expires and the day could end with a lower price. So these should be used with caution.",
      helpAskMarket: "A market sell is an order to sell equity shares at the current available market price. This order will be filled at the current market price, provided there is enough trading volume to process a trade. Market sells should be used when the investor wants to sell quickly. This could be because either the price is rapidly dropping and the investor is trying to cut losses, or to reinvest that money in a better stock.",
    }
    this.handleInfoPopup = this.handleInfoPopup.bind(this);
    this.handleCloseConfirmation = this.handleCloseConfirmation.bind(this);
  }

  handleInfoPopup(i) {
    // console.log(i);
    this.setState({
      help_type: i,
      showHelp: !this.state.showHelp,
    });
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

  inputChangeLimit(e) {
    var value =parseFloat(e.target.value);

    if (Number.isNaN(value)) {
      value = 0;
    }

    this.setState({
      limit_price: value
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
    var data_to_send ={"userIdString": JSON.parse(localStorage.getItem("profile")).sub, "equityTicker" : ticker, "amount" : this.state.number_of_stock, "limitPrice" : this.state.limit_price, "orderType" : this.state.action_type + this.props.button_name} ;
    var data = JSON.stringify(data_to_send);
    console.log(data);
    this.setState({renderSubmitted: true}); 
    jQuery.ajaxSetup({async:false});
    var url_type = this.props.button_name.toLowerCase();
    $.post(
      //"http://cloud-vm-45-112.doc.ic.ac.uk:8080/api/" + url_type,
      "http://localhost:8080/api/" + url_type,
      data,
      res => {
        this.setState({renderSubmitted: true})
      }
    );
  }

  render() {

    // Current amount is with respect to current price, that updates every graph poll. Can be changed if wanted. 

    var current_amount = this.props.current_price * this.state.number_of_stock;
    var amount_left = this.state.user_budget - current_amount;

    return (
      <div className="darken_bg">
        {this.state.renderSubmitted ? <SubmitConfirmation unmountMe={this.handleCloseConfirmation}/> : null}
        <div className="confirmation_window"> 
          <button className="close_button" onClick={() => this.dismiss()}> X </button> 
          <p className="company_viewing"> Viewing for {this.props.current_company} - {this.props.button_name}:</p>
          <div className="text_confirmation_window">
          <p> Number of stock: <input type="number" onChange={e => this.inputChangeStock(e)}/> </p>
          <p> Type of action: 
            {/* <div> */}
              <input type="radio" onClick={e => this.inputChangeAction(e)} id="actionChoice1" name="action" value="market" defaultChecked/>
              <label htmlFor="actionChoice1"> Market <a style={{color:"black"}} href="javascript:;" onClick={() => this.handleInfoPopup("market")}>(?)</a></label>

              <input type="radio" onClick={e => this.inputChangeAction(e)} id="actionChoice2" name="action" value="limit"/>
              <label htmlFor="actionChoice1"> Limit <a style={{color:"black"}} href="javascript:;" onClick={() => this.handleInfoPopup("limit")}>(?)</a></label>
              {this.state.action_type === "limit" ? <p> {this.props.button_name === "Bid" ? 'Maximum price to buy at' : 'Minimum price to sell at'}: <input type="number" onChange={e => this.inputChangeLimit(e)}/> </p> : null}
            {/* </div> */}
            <div className="temp_idea12">
              <div className="temp_idea13">
                {this.state.showHelp ? <InfoBubble2 text={
                  this.props.button_name === "Bid" ? (this.state.help_type === "limit" ? this.state.helpBidLimit : this.state.helpBidMarket) : (this.state.help_type === "limit" ? this.state.helpAskLimit : this.state.helpAskMarket) 
                }/> : null}
              </div>
            </div>
          </p>
          <p> Total price: {current_amount}</p>
          <p style={{color: amount_left > 0 ? "default" : "red"}}> Total funds left: {amount_left}</p> 
          <div className="place_order">
            <button className="place_order_button" disabled={amount_left < 0 || this.state.number_of_stock <= 0 || (this.state.action_type === 'limit' && this.state.limit_price <= 0)} onClick={() => this.submitRequest()}> Place order </button> 
          </div>
          </div>
        </div> 
      </div>
    )
  }
}

class InfoBubble2 extends React.Component {
  render() {
    return (
      <div className="info_bubble_wrapper">
        <div className="speech-bubble2"> 
          <p>
          {this.props.text} 
          </p>
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
