/** @jsx React.DOM */
'use strict'

var
	jQuery = require('jquery'),
	Page = require('../page/Page'),
	PaintingSummary = require('./PaintingSummary'),
	React = require('react'),
	underscore = require('underscore')
;

var PaintingsPage = React.createClass({
	getInitialState: function() {
		return {
			paintings: [],
			allowedWidth: window.innerWidth
		};
	},

	componentWillMount: function() {
		var me = this;
		jQuery.getJSON('/api/paintings', function(paintings) {
			var withPhotos = underscore.filter(paintings, function(painting) {
				return painting.Image.URL !== "";
			});
			me.setState({paintings: withPhotos});
		});
	},

	componentDidMount: function() {
    window.addEventListener('resize', this.handleResize);
  },

  componentWillUnmount: function() {
    window.removeEventListener('resize', this.handleResize);
  },

	render: function() {
		var me = this;
		return <Page>
			<div className='PaintingsPage' ref='paintings'>
				{(function() {
					var grouped = me.groupPaintings(me.state.paintings);

					return grouped.years.map(function(year) {
						return <section key={year}>
							<h1 ref={'year' + year}>{year}</h1>
							{(function() {
								var rows = me.justifyWidth(
									grouped.paintings[year],
									me.state.allowedWidth - 100,
									2
								);

								var i = 0;
								return rows.map(function(row) {
									i++;
									return <ul key={i}>{
										row.map(function(p) {
											var painting = p.painting;
											var dimensions = p.dimensions;
											return <li key={painting.ID}>
												<PaintingSummary painting={painting} dimensions={dimensions} />
											</li>;
										})
									}</ul>;
								});
							})()}
						</section>;
					})
				})()}
			</div>
		</Page>;
	},

	groupPaintings: function(paintings) {
		if (!paintings) {
			return {
				paintings: {},
				years: []
			};
		}

		var byYear = underscore.groupBy(this.state.paintings, function(painting) {
			return painting.Year;
		});
		var years = underscore.keys(byYear);
		years.sort();
		years.reverse();
		return {
			paintings: byYear,
			years: years
		};
	},

	justifyWidth: function(paintings, width, margin) {
		var me = this;
		var rows = [];
		var row = [];
		var accumulatedWidth = 0;
		var accumulatedMargin = 0;

		for (var i = 0; i < paintings.length; i++) {
			var painting = paintings[i];
			var dimensions = me.scaleToHeight(painting.Image, 300);
			row.push({painting: painting, dimensions: dimensions});
			accumulatedWidth += dimensions.width;
			accumulatedMargin += (margin || 0) * 2;

			if (width <= (accumulatedWidth + accumulatedMargin)) {
				// Rescale the entire row.
				var adjustedScale = (width - accumulatedMargin) / accumulatedWidth;
				row.forEach(function(p) {
					p.dimensions.width *= adjustedScale;
					p.dimensions.height *= adjustedScale;
				});

				rows.push(row);
				row = [];
				accumulatedWidth = 0;
				accumulatedMargin = 0;
			}
		}

		if (row.length > 0) {
			rows.push(row);
		}

		return rows;
	},

	scaleToHeight: function(image, height) {
		return {
			width: height * image.Width / image.Height,
			height: height 
		};
	},

	handleResize: underscore.debounce(function(event) {
		var width = this.refs['paintings'].getDOMNode().offsetWidth;
		this.setState({allowedWidth: width});
	}, 300)
});

module.exports = PaintingsPage;
