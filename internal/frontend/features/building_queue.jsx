import { h, render, Component } from 'preact';
import { route } from 'preact-router';
import axios from 'axios';
import classNames from 'classnames';
import arrayMove from 'array-move';

export default class BuildingQueue extends Component {
	state = {
		name: 'building queue',
		village_name: '',
		village_id: 0,
		all_villages: [],
		queue: [],
		error_village: false,
		buildings: [],
		resources: [],
		id_strict:0,
		buildings_dict: null
	}

	componentWillMount() {
		this.setState({
			...this.props.feature
		});
	}

	componentDidMount() {
		if (this.state.village_id) this.village_changes({ target: { value: this.state.village_id } });

		axios.get('/api/data?ident=villages').then(res => this.setState({ all_villages: res.data }));
		axios.get('/api/data?ident=buildingdata').then(res => this.setState({ buildings_dict: res.data }));
	}

	submit = async e => {
		this.setState({
			error_village: (this.state.village_id == 0)
		});

		if (this.state.error_village) return;

		const { ident, uuid, village_name, village_id, queue ,id_strict} = this.state;
		this.props.submit({ ident, uuid, village_name, village_id, queue,id_strict });
	}

	delete = async e => {
		const { ident, uuid, village_name, village_id, queue ,id_strict} = this.state;
		this.props.delete({ ident, uuid, village_name, village_id, queue,id_strict });
	}

	cancel = async e => {
		route('/');
	}

	village_changes = async (e) => {
		if (!e.target.value) return;
		this.setState({
			village_id: e.target.value
		});

		if(e.target[e.target.selectedIndex] != undefined){
			this.setState({
				village_name: e.target[e.target.selectedIndex].attributes.village_name.value
			});
		}
		let response = await axios.get(`/api/data?ident=buildings&village_id=${this.state.village_id}`);
		let res = [];
		let bd = [];
		const { buildings_dict , queue } = this.state;
		console.log(buildings_dict);
		let new_buildings_dict = [];//buildings_dict;
		for (let i = 1; i <=45;i++){
			new_buildings_dict[i] = JSON.parse(JSON.stringify(buildings_dict[i]));
			//new_buildings_dict[i]=buildings_dict[i];
		}
		console.log(new_buildings_dict);
		let new_item;
		for (let item of response.data) {
			if (Number(item.buildingType) > 4) {
				if (Number(item.lvl) > 0){
					new_buildings_dict[Number(item.buildingType)] = '-1';
				}
				new_item = JSON.parse(JSON.stringify(item));
				bd.push(item);
				continue;
			}

			res.push(item);
		}
		for (let i = 5; i<new_buildings_dict.length;i++){
			if (new_buildings_dict[i] != '-1'){
				new_item.buildingType = String(i);
				new_item.lvl = String(0);
				new_item.locationId = String(0);
				bd.push(JSON.parse(JSON.stringify(new_item)));
				continue;
			}
		}

		res = res.sort((x1, x2) => Number(x1.buildingType) - Number(x2.buildingType));
		bd = bd.sort((x1, x2) => Number(x1.buildingType) - Number(x2.buildingType));

		this.setState({
			buildings: bd,
			resources: res
		});
	}

	upgrade = building => {
		const { buildingType, lvl, locationId } = building;
		const { buildings_dict , queue, buildings } = this.state;
		let s_lvl = lvl;
		let loc  = [];
		for (let i = 0; i<=21; i++){loc[i]= String(i+19);}
		for (let i = 0; i<buildings.length;i++ ){
			if (Number(buildings[i].locationId) != 0){
				let index = loc.indexOf(String(buildings[i].locationId));
				if (index !== -1) loc.splice(index, 1);
			}
		}
		for (let i = 0; i<queue.length;i++ ){
			if (Number(queue[i].location) != 0){
				let index = loc.indexOf(String(queue[i].location));
				if (index !== -1) loc.splice(index, 1);
			}
		}
		if (queue != null) {
			for (let i=0; i<queue.length;i++){
				const { type, location } = queue[i];
				let index = loc.indexOf(String(locationId));
				if (index !== -1) loc.splice(index, 1);
				if (Number(type) == Number(buildingType) && Number(location) == Number(locationId)){
					s_lvl = Number(s_lvl)+1;
				}
			}
		}
		let c;
		let r;
		if (Number(s_lvl) > 0){
			let  resources1= [];
			resources1 = building.nextUpgradeCosts;
			let  resources2= [];
			resources2 = building.nextUpgradeTimes;
			r= resources2[Number(s_lvl)];
			if (r === null || r === undefined)
			{ r = building.upgradeTime;}
			c = resources1[Number(s_lvl)];
			if (c === null || c === undefined)
			{
				c = building.upgradeCosts;
			}
		} else {
			c = {
				'1' : 0,
				'2' : 0,
				'3' : 0,
				'4' : 0
			};
			r = 1;
		}
		let loct;
		let acc = '1';
		if (Number(s_lvl) == 0 && Number(locationId) > 0 ){
			if (Number(locationId) < 19 ){
				acc = '1';
			}else
			{
				acc = '0';
			}
			loc.push(String(locationId));
			s_lvl = '-1';

			c = {
				'1' : 0,
				'2' : 0,
				'3' : 0,
				'4' : 0
			};
			for (let i = 0; i<buildings.length;i++ ){
				if (Number(buildings[i].locationId) == Number(locationId) && Number(buildings[i].buildingType) == Number(buildingType) ){
					buildings.splice(i, 1);
					i = buildings.length+2;
				}
			}
		}
		console.log(locationId);
		if (Number(locationId) == 0){
			let is_find  = false;
			for(let i = 0; i < queue.length; i++)
			{
				if (Number(queue[i].actionBuild) > 0 && Number(queue[i].type) == Number(buildingType)) {
					acc = '1';
					loct = queue[i].location;
					if (Number(s_lvl) < Number(queue[i].build_lvl)) s_lvl = Number(queue[i].build_lvl);
					is_find = true;
					console.log(queue[is_find]);
				}
			}
			if(!is_find) {
				let ind = Math.floor(Math.random() * (loc.length - 1)) + 1;
				console.log(ind);
				acc = '2';
				console.log(loc);
				loct = loc[ind];
				console.log(loct);
				if (ind !== -1) loc.splice(ind, 1);
				console.log(loc);
			}
		} else {loct = locationId;}
		const queue_item = {
			actionBuild : acc,
			type: buildingType,
			build_name:  buildings_dict[Number(buildingType)],
			build_lvl: Number(s_lvl)+1,
			location: loct,
			costs: {
				//...building.upgradeCosts
				...c
			},
			//upgrade_time: building.upgradeTime
			upgrade_time: r
		};

		this.setState({ queue: [ ...this.state.queue, queue_item ], buildings: buildings });
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

	render({}, { name, all_villages, village_name, village_id, queue, buildings, resources, buildings_dict,id_strict }) {
		const village_select_class = classNames({
			select: true,
			'is-danger': this.state.error_village
		});

		const header_style = {
			textAlign: 'center'
		};
		const img_style = {
			width: 50,
			height: 50
		};

		let buildings_options = [];
		if (buildings_dict) {
			buildings_options = buildings.map(building =>
				<tr>
					<td style={ header_style }>{ building.locationId }</td>
					<td>{ buildings_dict[building.buildingType] }</td>
					<td style={ header_style }>{ building.lvl }</td>
					<td style={ header_style }>
						<a class="has-text-black" onClick={ e => this.upgrade(building) }>
							<span class="icon is-medium">
								<i class="far fa-lg fa-arrow-alt-circle-up"></i>
							</span>
						</a>
					</td>
				</tr>
			);
		}

		let resource_options = [];
		if (buildings_dict) {
			resource_options = resources.map(building =>
				<tr>
					<td style={ header_style }>{ building.locationId }</td>
					<td>{ buildings_dict[building.buildingType] }</td>
					<td style={ header_style }>{ building.lvl }</td>
					<td style={ header_style }>
						<a class="has-text-black" onClick={ e => this.upgrade(building) }>
							<span class="icon is-medium">
								<i class="far fa-lg fa-arrow-alt-circle-up"></i>
							</span>
						</a>
					</td>
				</tr>
			);
		}

		let queue_options = [];
		if (buildings_dict) {
			queue_options = queue.map((building, index) =>
				<tr>
					<td style={ header_style }>{ index + 1 }</td>
					<td style={ header_style }>{ building.location }</td>
					<td>{ buildings_dict[building.type] }</td>
					<td>{ building.build_lvl }</td>
					<td style={ header_style }>
						<a class="has-text-black" onClick={ e => this.move_up(building) }>
							<span class="icon is-medium">
								<i class="fas fa-lg fa-long-arrow-alt-up"></i>
							</span>
						</a>
					</td>
					<td style={ header_style }>
						<a class="has-text-black" onClick={ e => this.move_down(building) }>
							<span class="icon is-medium">
								<i class="fas fa-lg fa-long-arrow-alt-down"></i>
							</span>
						</a>
					</td>
					<td style={ header_style }>
						<a class="has-text-black" onClick={ e => this.delete_item(building) }>
							<span class="icon is-medium">
								<i class="far fa-lg fa-trash-alt"></i>
							</span>
						</a>
					</td>
				</tr>
			);
		}

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
										onload={this.village_changes}
										onChange={this.village_changes}

									>
										{ villages }
									</select>
								</div>
								<input type="checkbox"
									   style={ img_style }
									//checked="checked"
									//value={ units[unit.id] }
									   checked={ id_strict }
									   placeholder="0"

									   onChange={ async e => {
										   if (e.target.checked){
											   this.setState({ id_strict: 1 });
										   		id_strict   = 1;
										   } else
										   	{
												this.setState({ id_strict: 0 });
										   		id_strict   = 0;
										   	}
										   //this.setState({units_grate: units_grate});
										   //this.setfilter(e);
									   } }
								/>
								<a className="button is-success is-radiusless" style="margin-left: 3rem; margin-right: 1rem" onClick={ this.submit }>submit</a>
								<a className="button is-danger is-radiusless" onClick={ this.delete }>delete</a>

							</div>
						</div>
					</div>

				</div>

				<div className="columns" style="margin-top: 2rem">

					<div className="column" align="center">
						<strong>resource fields</strong>
						<table className="table is-striped">
							<thead>
								<tr>
									<td style={ header_style }><strong>id</strong></td>
									<td><strong>name</strong></td>
									<td style={ header_style }><strong>lvl</strong></td>
									<td style={ header_style }><strong></strong></td>
								</tr>
							</thead>
							<tbody>
								{ resource_options }
							</tbody>
						</table>
					</div>

					<div className="column" align="center">
						<strong>buildings</strong>
						<table className="table is-striped">
							<thead>
								<tr>
									<td style={ header_style }><strong>id</strong></td>
									<td><strong>name</strong></td>
									<td style={ header_style }><strong>lvl</strong></td>
									<td style={ header_style }><strong></strong></td>
								</tr>
							</thead>
							<tbody>
								{ buildings_options }
							</tbody>
						</table>
					</div>

					<div className="column" align="center">
						<strong>queue</strong>
						<table className="table is-striped">
							<thead>
								<tr>
									<td style={ header_style }><strong>pos</strong></td>
									<td style={ header_style }><strong>id</strong></td>
									<td><strong>name</strong></td>
									<td style={ header_style }><strong>lvl</strong></td>
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
