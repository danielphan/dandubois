/** @jsx React.DOM */
'use strict'

var
	jQuery = require('jquery'),
	React = require('react'),
  Router = require('react-router-component'),
  underscore = require('underscore')
;

var
	Link = Router.Link
;

underscore.once(function() {
	window.jQuery = jQuery;
	require('../../bower_components/unveil/jquery.unveil.js');
})();

var PaintingSummary = React.createClass({
	componentDidMount: function() {
		jQuery(this.refs['image'].getDOMNode()).unveil(100, function() {
			jQuery(this).load(function() {
		    this.style.opacity = 1;
		  });
		});
	},

	render: function() {
		var painting = this.props.painting;
		var dimensions = this.props.dimensions;

		// Take the longest side of the painting,
		// and make sure we can fetch an image of that size.
		var longest = 0;
		if (dimensions.width > dimensions.height) {
			longest = Math.min(painting.Image.Width, dimensions.width);
		} else {
			longest = Math.min(painting.Image.Height, dimensions.height);
		}
		// Max resize allowed by App Engine.
		if (longest > 1600) {
			longest = 0;
		}

		var longestRetina = longest * 2;
		if (longestRetina > 1600) {
			longestRetina = 0;
		}

		return <div className='PaintingSummary'>
			<Link href={'/paintings/' + painting.ID}>
				<p className='Title'>{painting.Title}</p>
				<img ref='image' className='Image'
					data-src={painting.Image.URL + '=s' + longest}
					data-src-retina={painting.Image.URL + '=s' + longestRetina}
					width={dimensions.width} height={dimensions.height} />
			</Link>
		</div>;
	}
});

module.exports = PaintingSummary;
