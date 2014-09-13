/** @jsx React.DOM */

var feedURL = 'https://medium.com/feed/latest'

/**
 * See https://github.com/dpup/strimmer for more details.
 */
var App = React.createClass({
  render: function () {
    var items = []
    return (
      <div>
        <Stream items={items} ref="stream" />
      </div>
    );
  },
  componentDidMount: function () {
    this.connect()
  },
  componentWillUnmount: function () {
    if (this.conn) this.conn.close()
    if (this.reconnectTimeout) clearTimeout(this.reconnectTimeout)
  },
  addItem: function (item) {
    this.refs.stream.addItem(item);
  },
  connect: function () {
    console.log('Connecting...')
    var conn = this.conn = new WebSocket('ws://' + location.host + '/bridge?feed=' + encodeURIComponent(feedURL));
    conn.onopen = this.onConnect;
    conn.onclose = this.onDisconnect;
    conn.onmessage = this.onRecieve;
  },
  onConnect: function (e) {
    console.log('Connection opened')
    this.addItem({
      Title: 'Connected to bridge...',
      Published: new Date().toISOString(),
    });
    this.connected = true;
  },
  onDisconnect: function (e) {
    console.log('Connection closed')
    if (this.connected) {
      // Only show if state changes.
      this.addItem({
        Title: 'Bridge disconnected',
        Published: new Date().toISOString(),
      });
    }
    this.reconnectTimeout = setTimeout(this.connect, 10000)
    this.conn = null
    this.connected = false;
  },
  onRecieve: function (e) {
    try {
      (JSON.parse(e.data).Entries || []).forEach(this.addItem);
    } catch (err) {
      console.error('Invalid message:', err, e.data);
    }
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
