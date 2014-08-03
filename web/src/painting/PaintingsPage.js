/** @jsx React.DOM */
'use strict'

var
	jQuery = require('jquery'),
	Page = require('../page/Page'),
	PaintingSummary = require('./PaintingSummary'),
	React = require('react'),
	_ = require('underscore')
;

var PaintingsPage = React.createClass({
	getInitialState: function() {
		return {
			paintings: [],
			width: window.innerWidth
		};
	},

	componentWillMount: function() {
		var me = this;
		jQuery.getJSON('/api/paintings', function(paintings) {
			var withPhotos = _(paintings).filter(function(painting) {
				return painting.Image.URL !== "";
			});
			me.setState({paintings: withPhotos});
		});
	},

	componentDidMount: function() {
		this.handleResize();
    window.addEventListener('resize', this.handleResize);
  },

  componentWillUnmount: function() {
    window.removeEventListener('resize', this.handleResize);
  },

  handleResize: _.debounce(function(event) {
		var width = this.refs['paintings'].getDOMNode().offsetWidth;
		this.setState({width: width});
	}, 300),

	render: function() {
		var me = this;
		return <Page>
			<div className='PaintingsPage' ref='paintings'>
				{(function() {
					var grouped = me.groupByYear(me.state.paintings);
					var years = _(grouped).keys();
					years.sort().reverse();
					return _(years).map(function(year) {
						var paintings = grouped[year];
						paintings = me.withDimensions(paintings);
						paintings.unshift({
							object: year,
							dimensions: {
								width: 300,
								height: 300
							}
						});
						var objects = me.justifyWidth(paintings, me.state.width - 100, 2);
						console.log(objects);
						return <ul key={year}>{
							_(objects).map(function(o, index) {
								if (index === 0) {
									// First element is the year.
									var style = {
										width: o.dimensions.width,
										height: o.dimensions.height,
										textAlign: 'center',
										fontSize: 80,
										paddingTop: 80
									};
									return <li key={'year-' + year} style={style}>{year}</li>;
								} else {
									var painting = o.object;
									return <li key={'painting-' + painting.ID}>
										<PaintingSummary painting={painting} dimensions={o.dimensions} />
									</li>;
								}
							})
						}</ul>;
					});
				})()}
			</div>
		</Page>;
	},

	groupByYear: function(paintings) {
		return _(paintings).groupBy(function(p) {
			return p.Year;
		});
	},

	withDimensions: function(paintings) {
		return _(paintings).map(function(p) {
			return {
				object: p,
				dimensions: {
					width: p.Width,
					height: p.Height
				}
			};
		});
	},

	// justifyWidth takse a list of `objects`:
	//
	//   object {
	//     dimensions: {
	//       width: number,
	//       height: number
	//     }
	//   }
	//
	// and mutates their dimensions so that they form rows such that each object
	// in a row has the same height, is separated by `margin` distance, and each
	// row's width is exactly `width`.
	justifyWidth: function(objects, width, margin) {
		var me = this;
		var rows = [];
		var row = [];
		var accumulatedWidth = 0;
		var accumulatedMargin = 0;

		objects.forEach(function(object) {
			object.dimensions = me.scaleToHeight(object.dimensions, 300);
			console.log(object);

			row.push(object);
			accumulatedWidth += object.dimensions.width;
			accumulatedMargin += (margin || 0) * 2;

			if (width <= (accumulatedWidth + accumulatedMargin)) {
				// Rescale the entire row.
				var adjustedScale = (width - accumulatedMargin) / accumulatedWidth;
				row.forEach(function(o) {
					o.dimensions.width *= adjustedScale;
					o.dimensions.height *= adjustedScale;
				});

				rows = rows.concat(row);
				row = [];
				accumulatedWidth = 0;
				accumulatedMargin = 0;
			}
		});

		if (row.length > 0) {
			rows = rows.concat(row);
		}

		rows.forEach(function(o) {
			o.dimensions.width = Math.floor(o.dimensions.width);
			o.dimensions.height = Math.floor(o.dimensions.height);
		});

		return rows;
	},

	scaleToHeight: function(dimensions, height) {
		return {
			width: dimensions.width * height / dimensions.height,
			height: height
		};
	}
});

module.exports = PaintingsPage;
