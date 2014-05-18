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

var Header = React.createClass({
  render: function() {
    return <header className='Header'>
    	<h1><span className='dan'>dan</span><span className='Dubois'>DuBois</span></h1>
      <nav>
      	<ul>
      		<li><Link href='/contact'>contact</Link></li>
      		<li><Link href='/login'>login</Link></li>
      	</ul>
      </nav>
    </header>;
  }
});

module.exports = Header;
