/** @jsx React.DOM */
'use strict'

var
	Page = require('../page/page'),
	React = require('react')
;

var ContactPage = React.createClass({
	render: function() {
		return <Page>
			<div className='ContactPage'>
				web@dandubois.net
			</div>
		</Page>;
	}
});

module.exports = ContactPage;
