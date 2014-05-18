/** @jsx React.DOM */
'use strict'

var
	Header = require('./Header'),
  React = require('react')
;

var Page = React.createClass({
  render: function() {
    return <div className='Page'>
    	<Header />
    	<div className='Content'>
    		{this.props.children}
    	</div>
    </div>;
  }
});

module.exports = Page;
