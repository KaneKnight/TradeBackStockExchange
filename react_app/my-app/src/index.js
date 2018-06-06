import React from 'react';
import ReactDOM from 'react-dom';
import $ from 'jquery';
import './stylesheets/style.css';

class App extends React.Component {

  constructor(props) {
    super(props);
    this.selectNewCompany = this.selectNewCompany.bind(this);
    this.state = {
      current_company: "Apple",
    };
  }

  selectNewCompany(new_company) {
    console.log("Called with " + new_company);
    console.log("Called");
    this.setState({current_company : new_company})
  }

  render() {
    //document.body.style.className = "c-container / t--light";
    return (
      //<div id='Stage' className="c-container / t--light">
      <div id='Stage' className="grid-container">
        <NavigationBar />
        <CompanyList onChange = {this.selectNewCompany}/>
        <GraphAndButtons current_company={this.state.current_company}/>
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
      opts: [],
    }
  }

  componentDidMount() {
    // Call code here to get company list
    //var source = ["Apple", "Apple 1", "Apple 2", "Apple 3", "Apple 4", "Apple 5", "Apple 6", "Apple 7", " Apple 8", "Apple 9"];
    var options = []
    //should be source.length
    for (var i = 0; i < 25; i++) {
      var source = "Apple " + i;
      options.push(<option key={source} value={source}> {source} </option>);
    }
    this.setState({
      opts: options,
    });

    // Code to set size for select. 45 magic number, don't change it please thanks.     
    var cont_div_height = $('#test1').height();
    var select_size = Math.floor(cont_div_height / 45);
    $('#company_sel').attr('size', select_size);
  }

  render() {
    return (
      <div id='test1' className="company_list_cont">
        <div className="list_of_companies">
          <select id='company_sel' className='company_select' onChange={(e) => this.props.onChange(e.target.value)}>
            {this.state.opts}
          </select>
        </div>
      </div>
    )
  }
}

class GraphAndButtons extends React.Component {
  render() {
    return (
      <div className="grid-container-graph"> 
        <Graph current_company={this.props.current_company}/>
        <UiInterface />
      </div>
    )
  }
}

class Graph extends React.Component {
  render() {
    return (
      <div className="graph_display_cont">
        <div className="graph_display"> Showing graph for {this.props.current_company}</div>
      </div>
    )
  }
}

function getFigures() {
  var temp = [2];
  temp[0] = 100;
  temp[1] = "+3.4";
  return temp;
}

class CompanyInfo extends React.Component {

  render() {

    var figures = getFigures();

    return (
      <div className="company_info_cont"> 
        Showing for {this.props.current_company}: 
        <br/> Price: {figures[0]}$
        <br/> Growth: {figures[1]}%
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
function Button(props) {
  return (
    <button className={props.button_type} onClick={props.onClick}> {props.button_name} </button>
  )
}

ReactDOM.render(<App />, document.getElementById('root'));
