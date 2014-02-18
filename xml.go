// Copyright 2014 Daniel Pupius

package main

import (
	"encoding/json"
	"encoding/xml"
	"log"
)

// Sample item.
//
// <?xml version="1.0" encoding="UTF-8" ?>
// <feed xmlns="http://www.w3.org/2005/Atom">
//     <status feed="https://medium.com/feed/latest" xmlns="http://superfeedr.com/xmpp-pubsub-ext">
//         <http code="200">Fetched (ping) 200 86400 and parsed 1/188 entries</http>
//         <next_fetch>2014-02-17T21:14:53.335Z</next_fetch>
//         <entries_count_since_last_maintenance>751</entries_count_since_last_maintenance>
//         <period>86400</period>
//         <last_fetch>2014-02-16T21:14:53.331Z</last_fetch>
//         <last_parse>2014-02-16T21:14:53.331Z</last_parse>
//         <last_maintenance_at>2014-02-15T10:15:22.000Z</last_maintenance_at>
//     </status>
//     <link title="Medium" rel="alternate" href="https://medium.com" type="text/html" />
//     <link title="Medium" rel="image" href="https://dnqgz544uhbo8.cloudfront.net/_/fp/img/default-preview-image.IsBK38jFAJBlWifMLO4z9g.png" type="image/png" />
//     <link title="Medium" rel="self" href="https://medium.com/feed/latest" type="application/rss+xml" />
//     <link title="" rel="hub" href="http://medium.superfeedr.com" type="text/html" />
//     <title>Medium</title>
//     <updated>2014-02-16T21:14:53.000Z</updated>
//     <id>medium-2014-2-16-21</id>
//     <entry xmlns="http://www.w3.org/2005/Atom" xmlns:geo="http://www.georss.org/georss" xmlns:as="http://activitystrea.ms/spec/1.0/" xmlns:sf="http://superfeedr.com/xmpp-pubsub-ext">
//         <id>https://medium.com/p/c9236332a756</id>
//         <published>2014-02-16T21:14:50.000Z</published>
//         <updated>2014-02-16T21:14:50.000Z</updated>
//         <title>Kissoda</title>
//         <summary type="html">&lt;div class="medium-feed-item"&gt;&lt;p class="medium-feed-snippet"&gt;Nothing; that&amp;#8217;s why I&amp;apos;m here;)&lt;/p&gt;&lt;/div&gt;</summary>
//         <link title="Kissoda" rel="alternate" href="https://medium.com/p/c9236332a756" type="text/html" />
//         <author>
//             <name>Kissoda</name>
//             <uri></uri>
//             <email></email>
//             <id>Kissoda</id>
//         </author>
//     </entry>
// </feed>

type Feed struct {
	Status  string  `xml:"status>http"`
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	URL         string `xml:"id"`
	Title       string `xml:"title"`
	Description string `xml:"summary"`
	Author      string `xml:"author>name"`
	Published   string `xml:"published"`
	Updated     string `xml:"updated"`
}

// XMLFeedToJSON unmarshals an RSS Feed using the Feed struct and then
// serializes it to JSON.
func XMLFeedToJSON(data []byte) ([]byte, error) {
	var feed Feed
	log.Println(string(data))
	err := xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	} else {
		return json.Marshal(feed)
	}
}
