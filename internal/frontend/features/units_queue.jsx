import { h, render, Component } from 'preact';
import { route } from 'preact-router';
import axios from 'axios';
import classNames from 'classnames';
import arrayMove from 'array-move';

export default class UnitsQueue extends Component {
	state = {
		name: 'units queue',
		village_name: '',
		village_id: 0,
		all_villages: [],
		queue: [],
		error_village: false,
		units: [],
		units_grate: [],
		units_dict: null
	}

	componentWillMount() {
		this.setState({
			...this.props.feature
		});
	}

	componentDidMount() {
		console.log(this.state.village_id);
		if (this.state.village_id) this.village_changes({ target: { value: this.state.village_id } });
		console.log(this.state.village_id);
		axios.get('/api/data?ident=villages').then(res => this.setState({ all_villages: res.data }));
		console.log(this.state.all_villages);
		axios.get('/api/data?ident=unitsdata').then(res => this.setState({ units_dict: res.data }));
		console.log(this.units_dict);
	}

	submit = async e => {
		this.setState({
			error_village: (this.state.village_id == 0)
		});

		if (this.state.error_village) return;

		const { ident, uuid, village_name, village_id, queue } = this.state;
		this.props.submit({ ident, uuid, village_name, village_id, queue });
	}

	delete = async e => {
		const { ident, uuid, village_name, village_id, queue } = this.state;
		this.props.delete({ ident, uuid, village_name, village_id, queue });
	}

	cancel = async e => {
		route('/');
	}
	Sec2Min(sec) {
		let min = (sec/60).toFixed(2);
		min = Number(Math.floor(min))+Number((((min%1)*60/100).toFixed(2))) ;
		return min;
	}
	village_changes = async e => {
		if (!e.target.value) return;
		this.setState({
			village_id: e.target.value
		});

		if(e.target[e.target.selectedIndex] != undefined){
			this.setState({
				village_name: e.target[e.target.selectedIndex].attributes.village_name.value
			});
		}
		let response_units = await axios.get(`/api/data?ident=units&is_great=true&village_id=${this.state.village_id}`);
		let response_units_great = await axios.get(`/api/data?ident=units&is_great=false&village_id=${this.state.village_id}`);
		let r_units = response_units.data.filter(element => element !== null);
		let r_units_great = response_units_great.data.filter(element => element !== null);

		this.setState({
			units: r_units,
			units_grate: r_units_great
		});
	}

	upgrade = unit => {
		if (unit.count != undefined && unit.run_time != undefined && Number(unit.count)>0 && Number(unit.run_time)>0)
		this.setState({ queue: [ ...this.state.queue, unit ] });
	}

	delete_item = building => {
		const queues = this.state.queue;
		var idx = queues.indexOf(building);
		if (idx != -1) {
			queues.splice(idx, 1); // The second parameter is the number of elements to remove.
		}

		this.setState({ queue: [ ...queues ] });
	}

	move_up = building => {
		const queues = this.state.queue;
		var idx = queues.indexOf(building);
		if (idx != -1) {
			arrayMove.mut(queues, idx, idx - 1); // The second parameter is the number of elements to remove.
		}

		this.setState({ queue: [ ...queues ] });
	}

	move_down = building => {
		const queues = this.state.queue;
		var idx = queues.indexOf(building);
		if (idx != -1) {
			arrayMove.mut(queues, idx, idx + 1); // The second parameter is the number of elements to remove.
		}

		this.setState({ queue: [ ...queues ] });
	}

	render({}, { name, all_villages, village_name, village_id, queue, units, units_grate, units_dict }) {
		const village_select_class = classNames({
			select: true,
			'is-danger': this.state.error_village
		});

		const header_style = {
			textAlign: 'center'
		};
		const img_style = {
			width: 93,
			height: 64
		};

		let units_options = [];
		console.log(units);
			units_options = units.map(unit =>
				<tr>
					<td> <img src= {'/images/'+unit.id+'.png'} alt= { unit.name }  style={ img_style }/></td>
					<td style={ header_style }>{ this.Sec2Min(unit.time) }</td>
					<td> <input
						style="width: 20px;"
						type="text"
						value={ unit.count }
						placeholder="0"
						onChange={ (e) => { unit.count =  e.target.value } }
					/></td>
					<td> <input
						style="width: 20px;"
						type="text"
						value={unit.run_time }
						placeholder="0"
						onChange={ (e) => { unit.run_time =  e.target.value } }
					/></td>
					<td style={ header_style }>
						<a class="has-text-black" onClick={ e => this.upgrade(unit) }>
							<span class="icon is-medium">
								<i class="far fa-lg fa-arrow-alt-circle-up"></i>
							</span>
						</a>
					</td>
				</tr>
			);


		let resource_options = [];
			resource_options = units_grate.map(unit =>
				<tr>
					<td> <img src= {'/images/'+unit.id+'.png'} alt= { unit.name }  style={ img_style }/></td>
					<td style={ header_style }>{ this.Sec2Min(unit.time) }</td>
					<td> <input
						style="width: 20px;"
						type="text"
						value={ unit.count }
						placeholder="0"
						onChange={ (e) => { unit.count =  e.target.value } }
					/></td>
					<td> <input
						style="width: 20px;"
						type="text"
						value={ unit.run_time }
						placeholder="0"
						onChange={ (e) => { unit.run_time =  e.target.value } }
					/></td>
					<td style={ header_style }>
						<a class="has-text-black" onClick={ e => this.upgrade(unit) }>
							<span class="icon is-medium">
								<i class="far fa-lg fa-arrow-alt-circle-up"></i>
							</span>
						</a>
					</td>
				</tr>
			);


		let queue_options = [];

			queue_options = queue.map((unit, index) =>
				<tr>
					<td> <img src= {'/images/'+unit.id+'.png'} alt= { unit.name }  style={ img_style }/></td>
					<td style={ header_style }>{ unit.count }</td>
					<td style={ header_style }>{ unit.run_time }</td>
					<td style={ header_style }>
						<a class="has-text-black" onClick={ e => this.move_up(unit) }>
							<span class="icon is-medium">
								<i class="fas fa-lg fa-long-arrow-alt-up"></i>
							</span>
						</a>
					</td>
					<td style={ header_style }>
						<a class="has-text-black" onClick={ e => this.move_down(unit) }>
							<span class="icon is-medium">
								<i class="fas fa-lg fa-long-arrow-alt-down"></i>
							</span>
						</a>
					</td>
					<td style={ header_style }>
						<a class="has-text-black" onClick={ e => this.delete_item(unit) }>
							<span class="icon is-medium">
								<i class="far fa-lg fa-trash-alt"></i>
							</span>
						</a>
					</td>
				</tr>
			);


		const villages = all_villages.map(village =>
			<option value={ village.data.villageId } village_name={ village.data.name } >({village.data.coordinates.x}|{village.data.coordinates.y}) {village.data.name}</option>
		);


		return (
			<div>
				<div className="columns">
					<div className="column">
						<div class="field">
							<label class="label">select village</label>
							<div class="control">
								<div class={ village_select_class }>
									<select
										class='is-radiusless'
										value={ village_id }
										onChange={ this.village_changes }
									>
										{ villages }
									</select>
								</div>
								<a className="button is-success is-radiusless" style="margin-left: 3rem; margin-right: 1rem" onClick={ this.submit }>submit</a>
								<a className="button is-danger is-radiusless" onClick={ this.delete }>delete</a>

							</div>
						</div>
					</div>
				</div>

				<div className="columns" style="margin-top: 2rem">

					<div className="column" align="center">
						<strong>Barracks, stable, workshop</strong>
						<table className="table is-striped">
							<thead>
							<tr>
								<td style={ img_style }><strong>\___Unit___/</strong></td>
								<td style={ header_style }><strong>time</strong></td>
								<td style={ header_style }><strong>count</strong></td>
								<td style={ header_style }><strong>period min</strong></td>
								<td style={ header_style }><strong></strong></td>
							</tr>
							</thead>
							<tbody>
							{ units_options }
							</tbody>
						</table>
					</div>

					<div className="column" align="center">
						<strong>great barracks, great stable</strong>
						<table className="table is-striped">
							<thead>
							<tr>
								<td style={ img_style }><strong>\___Unit___/</strong></td>
								<td style={ header_style }><strong>time</strong></td>
								<td style={ header_style }><strong>count</strong></td>
								<td style={ header_style }><strong>period min</strong></td>
								<td style={ header_style }><strong></strong></td>
							</tr>
							</thead>
							<tbody>
							{ resource_options }
							</tbody>
						</table>
					</div>

					<div className="column" align="center">
						<strong>queue</strong>
						<table className="table is-striped">
							<thead>
							<tr>
								<td style={ img_style }><strong>\___Unit___/</strong></td>
								<td style={ header_style }><strong>count</strong></td>
								<td style={ header_style }><strong>period min</strong></td>
								<td style={ header_style }><strong></strong></td>
								<td style={ header_style }><strong></strong></td>
								<td style={ header_style }><strong></strong></td>
								<td style={ header_style }><strong></strong></td>
							</tr>
							</thead>
							<tbody>
							{ queue_options }
							</tbody>
						</table>
					</div>
				</div>

			</div>
		);
	}
}
