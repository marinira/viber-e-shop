import { h, render, Component } from 'preact';
import axios from 'axios';
import classNames from 'classnames';
import Input from '../components/input';
import DoubleInput from '../components/double_input';
import { connect } from 'unistore/preact';
import { handle_response } from '../actions';
import InactiveTable from '../components/inactive_table';
import InfoTitle from '../components/info_title';
import { route } from 'preact-router';
//import { Ifarmlist_units } from '../../src/interfaces';


@connect('notifications', handle_response)
export default class InactiveFinder extends Component {
	state = {
		name: 'inactive finder',
		listName: '',
		listId: 0,
		maxEntriesCount:0,
		village_name: '',
		village_id: 0,
		all_farmlists: [],
		all_villages: [],
		error_village: false,
		error_farmlist: false,
		min_sec_farm: '',
		max_sec_farm: '',
		min_dist: '',
		max_dist: '',
		min_vill_farmlist: '',
		max_vill_farmlist: '',
		coef_unit: '',
		coef_time: '',
		villages_invectives: [],
		farmlists: [],
		loading: false,
		units: [],
		units_all: [],
		units_grate: [],
		units_dict: null,
		message: ''
	}

	description = 'searches for inactive players and displays their villages based on distance. \
		once you added them to your farmlist, you can use the easy scout feature to spy them.';

	componentWillMount() {



		this.setState({
			...this.props.feature
		});


	}
	componentDidMount() {
		axios.get('/api/data?ident=villages').then(res => {
			this.setState({ all_villages: res.data/*, village_id: res.data[0].villageId, village_name: res.data[0].data.name */});
		});

		axios.get('/api/data?ident=farmlists').then(res => this.setState({ all_farmlists: res.data }));

		axios.get('/api/data?ident=unitsdata').then(res => this.setState({ units_dict: res.data }));
		if (this.state.village_id > 0) {
			axios.get(`/api/data?ident=units&is_great=true&village_id=${this.state.village_id}`).then(res => this.setState({ units_all: res.data.filter(element => element !== null) }));

			axios.get(`/api/data?ident=units&is_great=false&village_id=${this.state.village_id}`).then(res => this.setState({ units_grate: res.data.filter(element => element !== null) }));
			console.log(this.state);
		}
	}

	clicked = async item => {
		const { selected_farmlist } = this.state;

		this.setState({
			error_farmlist: (selected_farmlist == '')
		});

		if (this.state.error_farmlist) return false;

		this.setState({ error_farmlist: false });

		const payload = {
			action: 'toggle',
			data: {
				farmlist: selected_farmlist,
				village: item
			}
		};

		let response = await axios.post('/api/inactivefinder', payload);
		console.log(response);
		const { error } = response.data;

		if (error) {
			this.props.handle_response(response.data);
			return false;
		}

		return true;
	}
	submit = async e => {
		this.setState({
			error_village: (this.state.village_id == 0)
		});

		if (this.state.error_village) return;

		const { ident, uuid, village_name, village_id ,	villages_invectives, farmlists ,listId, listName,maxEntriesCount, units, min_sec_farm, max_sec_farm, min_vill_farmlist, max_vill_farmlist,coef_time,coef_unit,min_dist,max_dist} = this.state;
		this.props.submit({ident, uuid, village_name, village_id, villages_invectives, farmlists, listId,listName,maxEntriesCount, units, min_sec_farm, max_sec_farm, min_vill_farmlist, max_vill_farmlist,coef_time,coef_unit,min_dist,max_dist });
	}

	delete = async e => {
		const { ident, uuid, village_name, village_id ,	villages_invectives, farmlists, listId, listName ,maxEntriesCount, units, min_sec_farm, max_sec_farm, min_vill_farmlist, max_vill_farmlist,coef_time,coef_unit,min_dist,max_dist} = this.state;
		this.props.delete({ident, uuid, village_name, village_id, villages_invectives, farmlists , listId, listName ,maxEntriesCount, units, min_sec_farm, max_sec_farm, min_vill_farmlist, max_vill_farmlist,coef_time,coef_unit,min_dist,max_dist });
	}

	cancel = async e => {
		route('/');
	}
	search = async e => {
		if (this.state.loading) return;

		this.setState({
			error_village: (this.state.village_id == 0)
		});

		if (this.state.error_village) return;

		this.setState({ loading: true, message: '', villages_invectives: [] });

		const {
			village_id
		} = this.state;

		let response = await axios.get(`/api/data?ident=inactive&village_id=${village_id}`);
		console.log(response)

		const { error, data, message } = response;//.data;
		console.log(response.data);
		console.log(data);

		this.setState({ villages_invectives: [ ...data ], loading: false });
		console.log('this.state.inactives');
		console.log(this.state.villages_invectives);
		if (error) {
			this.props.handle_response(response.data);
			return;
		}

		this.setState({ message });
	}

	tofarmlist = async e => {
		if (this.state.loading) return;

		this.setState({
			error_village: (this.state.village_id == 0)
		});

		if (this.state.error_village) return;

		const {
			maxEntriesCount,
			villages_invectives
		} = this.state;

		let i = 0;
		let list = 0;
		let newList = [];
		for (let v of villages_invectives){
			if (i >= maxEntriesCount){
				let list_id = {
					id: list,

				}

				i = 0;
				list++;
			}
		}
		console.log(response)

		const { error, data, message } = response;//.data;
		console.log(response.data);
		console.log(data);

		this.setState({ villages_invectives: [ ...data ], loading: false });
		console.log('this.state.villages_invectives');
		console.log(this.state.villages_invectives);
		if (error) {
			this.props.handle_response(response.data);
			return;
		}

		this.setState({ message });
	}

	changeVillage = async e => {
		axios.get(`/api/data?ident=units&is_great=true&village_id=${this.state.village_id}`).
		then(res => this.setState({units_all: res.data.filter(element => element !== null) }));

		axios.get(`/api/data?ident=units&is_great=false&village_id=${this.state.village_id}`).
		then(res => this.setState({units_grate: res.data.filter(element => element !== null) }));
		console.log("change village");
		console.log(this.state);

	}

	render({}, { name, villages_invectives, message,  min_dist,  max_dist, coef_unit,coef_time, loading, min_sec_farm, max_sec_farm, min_vill_farmlist, max_vill_farmlist }) {
		const { all_villages, all_farmlists, village_name, village_id, listId, listName, maxEntriesCount,
			units,
			units_all,
			units_grate,
			units_dict,} = this.state;

		const row_style = {
			verticalAlign: 'middle',
			textAlign: 'center',
		};


		const header_style = {
			textAlign: 'center'
		};
		const img_style = {
			width: 50,
			height: 50
		};
		const village_select_class = classNames({
			select: true,
			'is-danger': this.state.error_village
		});

		const farmlist_select_class = classNames({
			select: true,
			'is-danger': this.state.error_farmlist
		});

		const search_button = classNames({
			button: true,
			'is-success': true,
			'is-radiusless': true,
			'is-loading': loading
		});

		const villages = all_villages.map(village => <option value={ village.data.villageId } name={ village.data.name } >({village.data.coordinates.x}|{village.data.coordinates.y}) {village.data.name}</option>);
		console.log(villages);
		const farmlist_opt = all_farmlists.map(farmlist => <option value={ farmlist.data.listId }  listName={  farmlist.data.listName } maxEntriesCount = {farmlist.data.maxEntriesCount}  >{ farmlist.data.listName }</option>);
		console.log(farmlist_opt);

		let units_options = units_all.map(unit =>
			<tr>
			<td> <img src= {'/images/'+unit.id+'.png'} alt= { unit.name }  style={ img_style }/></td>
			<td style={ row_style }>
				<input type="checkbox"
					   style={ img_style }
					//checked="checked"
					//value={ units[unit.id] }
					   checked={ units[unit.id] }
					placeholder="0"

					onChange={ async e => {
						if (e.target.checked){units[unit.id]   = 1;} else {units[unit.id]   = 0;}
						 //this.setState({units_grate: units_grate});
						//this.setfilter(e);
					} }
				/>
			</td>
			</tr>
		);


		return (
			<div>
				<InfoTitle title={ name } description={ this.description } />

				<div className="columns">

					<div className="column">

						<div class="field">
							<label class="label">distance relative to</label>
							<div class="control">
								<div class={ village_select_class }>
									<select
										class="is-radiusless"
										value={ village_id }
										onChange={ (e) => {
											this.setState({
												village_name: e.target[e.target.selectedIndex].attributes.name.value,
												village_id: e.target.value
											});
											this.changeVillage(e);
										}
										}
									>
										{ villages }
									</select>
								</div>
							</div>
						</div>

						<DoubleInput
							label='Random sec (min / max)'
							placeholder1='default: 0'
							placeholder2='default: 500'
							value1={ min_sec_farm }
							value2={ max_sec_farm }
							onChange1={ e => this.setState({ min_sec_farm: e.target.value }) }
							onChange2={ e => this.setState({ max_sec_farm: e.target.value }) }
						/>

						<DoubleInput
							label='distance (min / max)'
							placeholder1='default: 0'
							placeholder2='default: 200'
							value1={ min_dist }
							value2={ max_dist }
							onChange1={ e => this.setState({ min_dist: e.target.value }) }
							onChange2={ e => this.setState({ max_dist: e.target.value }) }
						/>

						<button className={ search_button } onClick={ this.search } style='margin-right: 1rem'>
							search
						</button>
						<button className={ search_button } onClick={ this.tofarmlist } style='margin-right: 1rem'>
							To Fatm list
						</button>
						<div className="column">
							<button className="button is-radiusless is-success" onClick={ this.submit } style='margin-right: 1rem'>
								submit
							</button>
							<button className="button is-radiusless" onClick={ this.cancel } style='margin-right: 1rem'>
								cancel
							</button>

							<button className="button is-danger is-radiusless" onClick={ this.delete }>
								delete
							</button>
						</div>

					</div>
					<div className="column">

						<div class="field">
							<label class="label">add to farmlist</label>
							<div className="control">
								<div class={ farmlist_select_class }>
									<select
										class="is-radiusless"
										value={ listId }
										onChange={ e => this.setState({
											listName: e.target[e.target.selectedIndex].attributes.listName.value,
											maxEntriesCount: e.target[e.target.selectedIndex].attributes.maxEntriesCount.value,
												listId: e.target.value
											})
											}
									>
										{ farmlist_opt }
									</select>
								</div>
							</div>
						</div>

						<DoubleInput
							label='villages in farmlist (min / max)'
							placeholder1='default: 0'
							placeholder2='default: 100'
							value1={ min_vill_farmlist }
							value2={ max_vill_farmlist }
							onChange1={ e => this.setState({ min_vill_farmlist: e.target.value }) }
							onChange2={ e => this.setState({ max_vill_farmlist: e.target.value }) }
						/>

						<label class="label">coefficient units </label>
						<div class="field has-addons">
							<p class="control">
								<input
									class="input is-radiusless"
									type="text"
									placeholder="default: 5"
									value={ coef_unit }
									onChange={ e => this.setState({ coef_unit: e.target.value }) }
								/>
							</p>
						</div>
						<label className="label">coefficient Time farm </label>
						<div className="field has-addons">
							<p className="control">
								<input
									className="input is-radiusless"
									type="text"
									placeholder="default: 5"
									value={coef_time}
									onChange={e => this.setState({ coef_time: e.target.value })}
								/>
							</p>
						</div>

						<div className="content" style='margin-top: 1.5rem' >
							{ message }
						</div>

					</div>

				</div>
				<div>

					<table className="table is-hoverable is-fullwidth">
						<thead>
						<tr>
							{ units_options }
						</tr>

						</thead>
					</table>
				</div>

				<InactiveTable content={ villages_invectives } clicked={ this.clicked } />
			</div>
		);
	}
}
