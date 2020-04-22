import { h, render, Component } from 'preact';
import { route } from 'preact-router';
import axios from 'axios';
import classNames from 'classnames';

export default class FindNature extends Component {
	state = {
		name: 'find nature',
		own_tribe: 0,
		village_name: '',
		village_id: '',
		wait_time: '',
		all_villages: [],
		all_oasis: [],
		all_oasis_filter: [],
		target_x: '',
		target_y: '',
		target_villageId: '',
		target_village_name: '',
		target_playerId: '',
		target_player_name: '',
		target_tribeId: '',
		target_distance: '',
		troops: '',
		date: '',
		time: '',
		oas: [],
		u1: '',
		u2: '',
		u3: '',
		u4: '',
		u5: '',
		u6: '',
		u7: '',
		u8: '',
		u9: '',
		u10: '',
		u11: '',
		o1: '',
		o2: '',
		o3: '',
		o4: '',
		o5: '',
		o6: '',
		o7: '',
		o8: '',
		o9: '',
		o10: '',
		send_hero: false,
		error_wait_time: false,
		error_village: false,
		error_target: false
	}

	componentWillMount() {
		this.setState({
			...this.props.feature
		});
		for (let i = 1; i<11;i++)
			this.state.oas[i] = 0;
		axios.get('/api/data?ident=villages').then(res => this.setState({ all_villages: res.data }));
		axios.get('/api/data?ident=player_tribe').then(res => this.setState({ own_tribe: Number(res.data) }));
		axios.get('/api/data?ident=troops').then(res => this.setState({ troops: res.data }));
	}

	setTarget = async e => {
		this.setState({
			error_village: (this.state.village_id == 0)
		});

		if (this.state.error_village) {
			return;
		}

		const { village_id } = this.state;


		let response_natures = await axios.get(`/api/data?ident=natures&village_id=${village_id}`);
		console.log(response_natures);
		if (response_natures != undefined)
		{
			//let all_oasis = [];
			//for (let el of response_natures.data)
			//	all_oasis.push(el)
			this.setState({all_oasis: response_natures.data, all_oasis_filter: response_natures.data });
			//this.setfilter(e);

		}

		return;
	}

	submit = async e => {
		this.setState({
			error_wait_time: (this.state.wait_time == ''),
			error_village: (this.state.village_id == 0)
		});

		if (/*this.state.error_wait_time || */this.state.error_village) return;

		this.props.submit({ ...this.state });
	}

	delete = async e => {
		this.props.delete({ ...this.state });
	}

	cancel = async e => {
		route('/');
	}
	setfilter = async e => {

		this.state.oas[1] = this.state.o1;
		this.state.oas[2] = this.state.o2;
		this.state.oas[3] = this.state.o3;
		this.state.oas[4] = this.state.o4;
		this.state.oas[5] = this.state.o5;
		this.state.oas[6] = this.state.o6;
		this.state.oas[7] = this.state.o7;
		this.state.oas[8] = this.state.o8;
		this.state.oas[9] = this.state.o9;
		this.state.oas[10] = this.state.o10;
		let count = 0;
		for (let i = 1; i<11; i++){
			count = count+  this.state.oas[i];
		}
		if (count !=0) {
			let fill = [];
			for (let oasis of this.state.all_oasis) {
				let is_exclude = true;
				for (let i = 1; i < 11; i++) {
					if ((oasis.units[i].count == 0 && this.state.oas[i] == 1)||
						(oasis.units[i].count > 0  && this.state.oas[i] == 0)||
						(oasis.units[i].count > 0 && this.state.oas[i] == -1)) {
							is_exclude = false;
					}
				}
				if (is_exclude ) {
					fill.push(oasis);
				}
			}
			this.setState({ all_oasis_filter: fill });
		} else {
			this.setState({ all_oasis_filter: this.state.all_oasis });
		}
	}

	render() {
		var { wait_time, all_villages, village_name, village_id, target_x, target_y, target_player_name, target_village_name, target_tribeId, target_distance, own_tribe, troops,
			u1,
			u2,
			u3,
			u4,
			u5,
			u6,
			u7,
			u8,
			u9,
			u10,
			t11,
			o1,
			o2,
			o3,
			o4,
			o5,
			o6,
			o7,
			o8,
			o9,
			o10,
			send_hero,
			time,
			date,
			all_oasis,
			all_oasis_filter
		} = this.state;

		var new_rows = [];
		if (own_tribe != 0 && troops != '') {
			new_rows = [
				<th style={ row_style }> {troops[own_tribe][1].name} </th>,
				<th style={ row_style }> {troops[own_tribe][2].name} </th>,
				<th style={ row_style }> {troops[own_tribe][3].name} </th>,
				<th style={ row_style }> {troops[own_tribe][4].name} </th>,
				<th style={ row_style }> {troops[own_tribe][5].name} </th>,
				<th style={ row_style }> {troops[own_tribe][6].name} </th>,
				<th style={ row_style }> {troops[own_tribe][7].name} </th>,
				<th style={ row_style }> {troops[own_tribe][8].name} </th>,
				<th style={ row_style }> {troops[own_tribe][9].name} </th>,
				<th style={ row_style }> {troops[own_tribe][10].name} </th>,
			];
		}


		if (date == '') {
			var curDate = new Date();
			curDate = curDate.toJSON();
			date = curDate.split('T')[0];
			this.setState({ date });
		}
		if (this.state.time == '') {
			var curUTCTime = new Date();
			curUTCTime = curUTCTime.toJSON();
			time = curUTCTime.split('T')[1].substring(0, 5);
			this.setState({ time });
		}

		const input_wait_time = classNames({
			input: true,
			'is-radiusless': true,
			'is-danger': this.state.error_wait_time
		});

		const village_select_class = classNames({
			select: true,
			'is-danger': this.state.error_village
		});

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

		const villages = all_villages.map(village => <option value={ village.data.villageId } village_name={ village.data.name } >({village.data.coordinates.x}|{village.data.coordinates.y}) {village.data.name}</option>);

		let queue_options = [];
	console.log(all_oasis);
		if (all_oasis_filter.length > 0) {
			queue_options = all_oasis_filter.map((oasis, index) =>

				<tr>
					<td style={header_style}> ({oasis.x}|{oasis.y})</td>
					<td style={ img_style }> <img src= {'/images/o'+oasis.oasisType+'.png'} alt= 'type' style={ img_style }/></td>
					<td style={header_style}> ({oasis.units[1].count})</td>
					<td style={header_style}> ({oasis.units[2].count})</td>
					<td style={header_style}> ({oasis.units[3].count})</td>
					<td style={header_style}> ({oasis.units[4].count})</td>
					<td style={header_style}> ({oasis.units[5].count})</td>
					<td style={header_style}> ({oasis.units[6].count})</td>
					<td style={header_style}> ({oasis.units[7].count})</td>
					<td style={header_style}> ({oasis.units[8].count})</td>
					<td style={header_style}> ({oasis.units[9].count})</td>
					<td style={header_style}> ({oasis.units[10].count})</td>
					<td style={header_style}>{oasis.dist}</td>
					<td style={header_style}>
						<a class="has-text-black" onClick={e => this.delete_item(oasis)}>
							<span class="icon is-medium">
								<i class="far fa-lg fa-trash-alt"></i>
							</span>
						</a>
					</td>
				</tr>

			);
		}

		return (
			<div>
				<div className="columns">

					<div className="column">
						<div>
							<label class="label">Target Land Time: UTC</label>
							<input type="date" id="start" name="trip-start"
								   value={ date } onChange={ (e) => this.setState({ date: e.target.value }) }
							></input>
							<input type="time" id="meeting-time" step="1"
								   name="meeting-time" value={ time } onChange={ (e) => this.setState({ time: e.target.value }) }
							/>
						</div>
						<div>
							<label class="label">x</label>
							<input
								style="width: 150px;"
								type="text"
								value={ target_x }
								placeholder="0"
								onChange={ (e) => this.setState({ target_x: e.target.value }) }
							/>
							<label class="label">y</label>
							<input
								style="width: 150px;"
								type="text"
								value={ target_y }
								placeholder="0"
								onChange={ (e) => this.setState({ target_y: e.target.value }) }
							/>

							<label class="label">send hero</label>
							<input type="checkbox" value={ send_hero } onChange={ (e) => this.setState({ send_hero: e.target.checked }) } />
						</div>



							<button className='button is-radiusless is-success' style='margin-top: 1rem' onClick={ this.setTarget }>
								set target
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
							<label class="label">select village</label>
							<div class="control">
								<div class={ village_select_class }>
									<select
										class="is-radiusless"
										value={ village_id }
										onChange={ (e) => this.setState({
											village_name: e.target[e.target.selectedIndex].attributes.village_name.value,
											village_id: e.target.value
										})
										}
									>
										{villages}
									</select>
								</div>
							</div>
						</div>


					</div>

				</div>

				<div>
					<table className="table is-hoverable is-fullwidth">
						<thead>
						<tr>
							<th style={ row_style }>distance</th>
							<th style={ row_style }>player</th>
							<th style={ row_style }>village</th>
							{new_rows}
						</tr>
						</thead>
						<tbody>
						<tr>
							<td style={ row_style }>
								{target_distance}
							</td>
							<td style={ row_style }>
								{target_player_name}
							</td>
							<td style={ row_style }>
								{target_village_name}
							</td>
							<td style={ row_style }>
								<input
									style="width: 30px;"
									type="text"
									value={ u1 }
									placeholder="u1"
									onChange={ async e => {
										this.setState({ u1: e.target.value });
									} }
								/>
							</td>
							<td style={ row_style }>
								<input
									style="width: 30px;"
									type="text"
									value={ u2 }
									placeholder="u2"
									onChange={ async e => {
										this.setState({ u2: e.target.value });
									} }
								/>
							</td>
							<td style={ row_style }>
								<input
									style="width: 30px;"
									type="text"
									value={ u3 }
									placeholder="u3"
									onChange={ async e => {
										this.setState({ u3: e.target.value });
									} }
								/>
							</td>
							<td style={ row_style }>
								<input
									style="width: 30px;"
									type="text"
									value={ u4 }
									placeholder="u4"
									onChange={ async e => {
										this.setState({ u4: e.target.value });
									} }
								/>
							</td>
							<td style={ row_style }>
								<input
									style="width: 30px;"
									type="text"
									value={ u5 }
									placeholder="u5"
									onChange={ async e => {
										this.setState({ u5: e.target.value });
									} }
								/>
							</td>
							<td style={ row_style }>
								<input
									style="width: 30px;"
									type="text"
									value={ u6 }
									placeholder="u6"
									onChange={ async e => {
										this.setState({ u6: e.target.value });
									} }
								/>
							</td>
							<td style={ row_style }>
								<input
									style="width: 30px;"
									type="text"
									value={ u7 }
									placeholder="u7"
									onChange={ async e => {
										this.setState({ u7: e.target.value });
									} }
								/>
							</td>
							<td style={ row_style }>
								<input
									style="width: 30px;"
									type="text"
									value={ u8 }
									placeholder="u8"
									onChange={ async e => {
										this.setState({ u8: e.target.value });
									} }
								/>
							</td>
							<td style={ row_style }>
								<input
									style="width: 30px;"
									type="text"
									value={ u9 }
									placeholder="u9"
									onChange={ async e => {
										this.setState({ u9: e.target.value });
									} }
								/>
							</td>
							<td style={ row_style }>
								<input
									style="width: 30px;"
									type="text"
									value={ u10 }
									placeholder="u10"
									onChange={ async e => {
										this.setState({ u10: e.target.value });
									} }
								/>
							</td>
						</tr>
						</tbody>
					</table>
				</div>


					<div>

						<table className="table is-hoverable is-fullwidth">
							<thead>
							<tr>
								<td style={ header_style }> (x|y)</td>
								<td style={ header_style }> type oasis </td>
								<td> <img src= {'/images/u1.png'} alt= 'Rat'  style={ img_style }/></td>
								<td> <img src= {'/images/u2.png'} alt= 'Spider' style={ img_style }/></td>
								<td> <img src= {'/images/u3.png'} alt= 'Serpant'  style={ img_style }/></td>
								<td> <img src= {'/images/u4.png'} alt= 'Bat'  style={ img_style }/></td>
								<td> <img src= {'/images/u5.png'} alt= 'Wild Boar'  style={ img_style }/></td>
								<td> <img src= {'/images/u6.png'} alt= 'Wolf' style={ img_style }/></td>
								<td> <img src= {'/images/u7.png'} alt= 'Bear'  style={ img_style }/></td>
								<td> <img src= {'/images/u8.png'} alt= 'Crocodile'  style={ img_style }/></td>
								<td> <img src= {'/images/u9.png'} alt= 'Tiger'  style={ img_style }/></td>
								<td> <img src= {'/images/u10.png'} alt= 'Elephant' style={ img_style }/></td>
								<td style={ header_style }>Distans</td>
							</tr>
							<tr>
								<td style={ row_style }>

								</td>
								<td> <img src= {'/images/o10.png'} alt= 'type' style={ img_style }/></td>
								<td style={ row_style }>
									<input
										style="width: 30px;"
										type="text"
										value={ o1 }
										placeholder="o1"
										onChange={ async e => {
											this.setState({ o1: e.target.value });
											this.setfilter(e);
										} }
									/>
								</td>
								<td style={ row_style }>
									<input
										style="width: 30px;"
										type="text"
										value={ o2 }
										placeholder="o2"
										onChange={ async e => {
											this.setState({ o2: e.target.value });
											this.setfilter(e);
										} }
									/>
								</td>
								<td style={ row_style }>
									<input
										style="width: 30px;"
										type="text"
										value={ o3 }
										placeholder="o3"
										onChange={ async e => {
											this.setState({ o3: e.target.value });
											this.setfilter(e);
										} }
									/>
								</td>
								<td style={ row_style }>
									<input
										style="width: 30px;"
										type="text"
										value={ o4 }
										placeholder="o4"
										onChange={ async e => {
											this.setState({ o4: e.target.value });
											this.setfilter(e);
										} }
									/>
								</td>
								<td style={ row_style }>
									<input
										style="width: 30px;"
										type="text"
										value={ o5 }
										placeholder="o5"
										onChange={ async e => {
											this.setState({ o5: e.target.value });
											this.setfilter(e);
										} }
									/>
								</td>
								<td style={ row_style }>
									<input
										style="width: 30px;"
										type="text"
										value={ o6 }
										placeholder="o6"
										onChange={ async e => {
											this.setState({ o6: e.target.value });
											this.setfilter(e);
										} }
									/>
								</td>
								<td style={ row_style }>
									<input
										style="width: 30px;"
										type="text"
										value={ o7 }
										placeholder="o7"
										onChange={ async e => {
											this.setState({ o7: e.target.value });
											this.setfilter(e);
										} }
									/>
								</td>
								<td style={ row_style }>
									<input
										style="width: 30px;"
										type="text"
										value={ o8 }
										placeholder="o8"
										onChange={ async e => {
											this.setState({ o8: e.target.value });
											this.setfilter(e);
										} }
									/>
								</td>
								<td style={ row_style }>
									<input
										style="width: 30px;"
										type="text"
										value={ o9 }
										placeholder="o9"
										onChange={ async e => {
											this.setState({ o9: e.target.value });
											this.setfilter(e);
										} }
									/>
								</td>
								<td style={ row_style }>
									<input
										style="width: 30px;"
										type="text"
										value={ o10 }
										placeholder="o10"
										onChange={ async e => {
											this.setState({ o10: e.target.value });
											this.setfilter(e);
										} }
									/>
								</td>
							</tr>
							</thead>
							<tbody>
							{ queue_options }
							</tbody>
						</table>
					</div>


			</div>
		);
	}
}
