/** @jsx React.DOM */
'use strict'

var
	Page = require('../page/Page'),
	React = require('react'),
	underscore = require('underscore')
;

var PaintingPage = React.createClass({
	getInitialState: function() {
		return {
			width: window.innerWidth
		};
	},

	componentWillMount: function() {
		var me = this;
		jQuery.getJSON('/api/paintings/' + this.props.id, function(painting) {
			me.setState({painting: painting});
		});
	},

	componentDidMount: function() {
		this.handleResize();
    window.addEventListener('resize', this.handleResize);
  },

  componentWillUnmount: function() {
    window.removeEventListener('resize', this.handleResize);
  },

  handleResize: underscore.debounce(function(event) {
		var width = this.refs['painting'].getDOMNode().offsetWidth;
		this.setState({width: width});
	}, 300),

	render: function() {
		var me = this;
		return <Page>
			<div className='PaintingPage' ref='painting'>
				{(function() {
					var painting = me.state.painting;
					if (!painting) {
						return;
					}

					var dimensions = {
						width: painting.Image.Width,
						height: painting.Image.Height
					};

					if (me.state.width < dimensions.width) {
						var scale = me.state.width / dimensions.width;
						dimensions.width = Math.floor(dimensions.width * scale);
						dimensions.height = Math.floor(dimensions.height * scale);
					}

					var longest = 0;
					if (dimensions.width > dimensions.height) {
						longest = me.state.width;
					}
					if (longest > 1600) {
						longest = 0;
					}

					return [
						<h1 className='Title'>{painting.Title}</h1>,
						<img ref='image' className='Image'
							src={painting.Image.URL + '=s' + longest}
							width={dimensions.width} height={dimensions.height} />
					];
				})()}
			</div>
		</Page>;
	}
});

module.exports = PaintingPage;
