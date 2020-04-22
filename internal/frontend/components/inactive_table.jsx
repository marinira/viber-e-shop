import { h, render, Component } from 'preact';
import { route } from 'preact-router';
import classNames from 'classnames';

export default class InactiveTable extends Component {
	render({ content, clicked }) {
		const row_style = {
			verticalAlign: 'middle',
			textAlign: 'center'
		};

		const list = content.map(item => <Inactive content={ item } clicked={ clicked } />);

		return (
			<div>
				<table className="table is-hoverable is-fullwidth">
					<thead>
						<tr>
							<th style={ row_style }>distance</th>
							<th style={ row_style }>population</th>
							<th style={ row_style }>coordinates</th>
							<th style={ row_style }>player</th>
							<th style={ row_style }>village</th>
							<th style={ row_style }>id villag</th>
							<th style={ row_style }>farmlist</th>
							<th />
						</tr>
					</thead>
					<tbody>
						{ list }
					</tbody>
				</table>
			</div>
		);
	}
}

class Inactive extends Component {
	state = {
		toggled: false
	}



	render({ content, clicked }, { toggled }) {

		let { dist, playerId, name, population, x, y, id, farmlist } = content;
		dist = Number(dist).toFixed(1);
		const coordinates = `( ${x} | ${y} )`;
		const tribe = farmlist;

		const row_style = {
			verticalAlign: 'middle',
			textAlign: 'center'
		};

		const icon = classNames({
			'fas': true,
			'fa-lg': true,
			'fa-plus': !toggled,
			'fa-minus': toggled
		});

		return (
			<tr>
				<td style={ row_style }>
					{ dist }
				</td>
				<td style={ row_style }>
					{ population }
				</td>
				<td style={ row_style }>
					{ coordinates }
				</td>
				<td style={ row_style }>
					{ playerId }
				</td>
				<td style={ row_style }>
					{ name }
				</td>
				<td style={ row_style }>
					{ id }
				</td>
				<td style={ row_style }>
					{ farmlist }
				</td>
				<td style={ row_style }>
					<a class="has-text-black" onClick={ async e => {
						if (await clicked(content)) this.setState({ toggled: !toggled });
					} }>
						<span class="icon is-medium">
							<i class={ icon }></i>
						</span>
					</a>
				</td>
			</tr>
		);
	}
}
