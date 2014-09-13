/** @jsx React.DOM */

/**
 * Provides a very simple WebSocket client, that can be used to connect to the
 * bridge. The default feed is from http://push-pub.appspot.com and allows you
 * to send test items.  The bridge location should be a publicly accessible host
 * or IP that routes to your bridge.
 *
 * See https://github.com/dpup/strimmer for more details.
 */
var App = React.createClass({
  render: function () {
    var items = []
    return (
      <div>
        <Form onSubscribe={this.onSubscribe} />
        <Stream items={items} ref="stream" />
      </div>
    );
  },
  onSubscribe: function (state) {
    var conn = new WebSocket('ws://' + state.bridgeURL + '/bridge?feed=' + encodeURIComponent(state.feedURL));
    conn.onclose = function (e) {
      this.addItem({
        Title: 'Disconnected from ' + state.feedURL,
        Published: new Date().toISOString(),
      });
    }.bind(this);
    conn.onmessage = function (e) {
      try {
        (JSON.parse(e.data).Entries || []).forEach(this.addItem);
      } catch (err) {
        console.error('Invalid message:', err, e.data);
      }
    }.bind(this);
    this.addItem({
      Title: 'Connected to ' + state.feedURL,
      Published: new Date().toISOString(),
    });
  },
  addItem: function (item) {
    this.refs.stream.addItem(item);
  }
});


/**
 * Form is a component that handles user input for subscribing to a new feed.
 */
var Form = React.createClass({
  mixins: [React.addons.LinkedStateMixin],
  getInitialState: function () {
    return {
      bridgeURL: location.host,
      feedURL: 'http://push-pub.appspot.com/feed'
    };
  },
  render: function () {
    return (
      <div className="form">
        <input valueLink={this.linkState('bridgeURL')} placeholder="my.bridge.com:1234" /> Bridge location<br />
        <input valueLink={this.linkState('feedURL')} /> Feed URL<br />
        <button onClick={this.onSubscribeClicked}>Subscribe</button>
      </div>
    );
  },
  onSubscribeClicked: function (e) {
    this.props.onSubscribe(this.state)
  }
});


/**
 * Stream is a component that renders a list of feed items.
 */
var Stream = React.createClass({
  addItem: function (item) {
    item.Id = item.URL || Math.random();
    var exists = this.state.items.some(function (o) {
      return o.Id === item.Id
    })
    if (!exists) {
      // TOOD(dan): This seems a bit of an ugly way to do this.
      var items = [item].concat(this.state.items);
      items.length = Math.min(items.length, 50);
      this.setState({items: items});
    } else {
      console.log('Duplicate item recieved', item.URL)
    }
  },
  getInitialState: function () {
    return {
      now: Date.now(),
      items: []
    };
  },
  componentDidMount: function() {
    this.interval = setInterval(this.tick, 1000);
  },
  componentWillUnmount: function() {
    clearInterval(this.interval);
  },
  render: function () {
    var now = this.state.now;
    var items = this.state.items;
    return (
      <div className="canvas">
        {items.map(function (item) {
          if (item.URL) {
            // Feed item.
            return (
              <article key={item.Id}>
                <h2><a href={item.URL} target="_blank">{item.Title}</a></h2>
                <div className="content">
                  <span dangerouslySetInnerHTML={{__html: item.Description}} />
                  <p className="info">
                    <time dateTime={item.Published}>
                      {this.relativeDate(item.Published, now)}
                    </time> by {item.Author}
                  </p>
                </div>
              </article>
            );
          } else {
            // System message.
            return (
              <article key={item.Id}>
                <p className="info">{item.Title}</p>
              </article>
            );
          }
        }.bind(this))}
      </div>
    );
  },
  tick: function () {
    this.setState({now: Date.now()});
  },
  relativeDate: function (dateString, now) {
    var diff = now - new Date(dateString);
    diff /= 1000;
    if (diff < 60) return this.relativeString(diff, 'second', 5);
    diff /= 60;
    if (diff < 60) return this.relativeString(diff, 'minute', 1);
    diff /= 60;
    if (diff < 48) return this.relativeString(diff, 'hour', 1);
    return 'some time ago';
  },
  relativeString: function (count, unit, roundTo) {
    count = Math.round(count / roundTo) * roundTo;
    return count + ' ' + unit + (count == 1 ? '' : 's') + ' ago';
  }
})


React.renderComponent(
  <App />,
  document.body
);
