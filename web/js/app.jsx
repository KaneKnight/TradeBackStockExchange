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

  generateDummyOptions() {
    var result = [];
    var dummy_data_str = {"packet" : "hi"};
    var dummy_data = JSON.stringify(dummy_data_str);
    jQuery.ajaxSetup({async:false});
    $.get(
      "http://localhost:8080/api/get-company-list",
      res => {
       // console.log(res.results);
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
  render() {
    return (
      <div className="graph_display_cont">
        <div className="graph_display"> Showing graph for {this.props.current_company}</div>
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

  render() {

    //var figures = getFigures(this.props.current_company);
    var figures = [1, 2];
    console.log("updated");

    return (
      <div className="company_info_cont"> 
        Showing for {this.props.current_company}: 
        <br/> Price: {figures[0]}$
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
    var dummy_data_buy = {"userId" : 1, "equityTicker" : this.props.current_company, "amount" : 1, "orderType" : "marketBid"};
    this.serverRequest(dummy_data_buy, "bid");
  }

  sell() {
    var dummy_data_sell = {"userId" : 1, "equityTicker" : this.props.current_company, "amount" : 1, "orderType" : "marketAsk"};
    this.serverRequest(dummy_data_sell, "ask");
  }

  serverRequest(dummy_data_str, url_type) {
    var dummy_data = JSON.stringify(dummy_data_str);
    /* TODO: Change local host to the actual address of the server. */
    console.log("Sent POST request for request:" + url_type);
    $.post(
      //"http://cloud-vm-45-112.doc.ic.ac.uk:8080/api/" + url_type,
      "http://localhost:8080/api/" + url_type,
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
function Button(props) {
  return (
    <button className={props.button_type} onClick={props.onClick}> {props.button_name} </button>
  )
}

ReactDOM.render(<App />, document.getElementById('app'));
