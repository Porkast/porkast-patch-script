package search

const CREATE_FEED_ITEM_INDEX_REQUEST = `
{
    "name": "feed_item",
    "storage_type": "disk",
    "shard_num": 1,
    "mappings": {
        "properties": {
            "title": {
                "type": "text",
                "index": true,
                "store": true,
                "highlightable": true,
                "analyzer": "gse_search",
                "search_analyzer": "gse_standard"
            },
            "author": {
                "type": "text",
                "index": true,
                "store": true,
                "highlightable": true,
                "analyzer": "gse_search",
                "search_analyzer": "gse_standard"
            },
            "description": {
                "type": "text",
                "index": true,
                "store": true,
                "highlightable": true,
                "analyzer": "gse_search",
                "search_analyzer": "gse_standard"
            },
            "textDescription": {
                "type": "text",
                "index": true,
                "store": true,
                "highlightable": true,
                "analyzer": "gse_search",
                "search_analyzer": "gse_standard"
            },
            "pubDate": {
                "type": "date",
                "index": true,
                "format": "2006-01-02 15:04:05",
                "sortable": true,
                "aggregatable": true
            },
            "created": {
                "type": "date",
                "store": true,
                "index": false
            },
            "id": {
                "type": "text",
                "store": true,
                "index": false
            },
            "channelId": {
                "type": "text",
                "store": true,
                "index": false
            },
            "link": {
                "type": "text",
                "store": true,
                "index": false
            },
            "imageUrl": {
                "type": "text",
                "store": true,
                "index": false
            },
            "enclosureUrl": {
                "type": "text",
                "store": true,
                "index": false
            },
            "enclosureType": {
                "type": "text",
                "store": true,
                "index": false
            },
            "enclosureLength": {
                "type": "text",
                "store": true,
                "index": false
            },
            "duration": {
                "type": "keyword",
                "index": true,
                "store": true
            },
            "episode": {
                "type": "text",
                "index": true,
                "store": true
            },
            "explicit": {
                "type": "text",
                "store": true,
                "index": false
            },
            "season": {
                "type": "text",
                "index": false,
                "store": true
            },
            "episodeType": {
                "type": "text",
                "store": true,
                "index": true
            },
            "channelImageUrl": {
                "type": "text",
                "store": true,
                "index": false
            },
            "channelTitle": {
                "type": "text",
                "index": true,
                "store": true,
                "analyzer": "gse_search",
                "search_analyzer": "gse_standard"
            },
            "language": {
                "type": "keyword",
                "index": true,
                "store": true,
                "analyzer": "gse_search",
                "search_analyzer": "gse_standard"
            },
            "feedLink": {
                "type": "text",
                "store": true,
                "index": false
            }
        }
    }
}
`

const CREATE_FEED_CHANNEL_INDEX_REQUEST = `
{
    "name": "feed_channel",
    "storage_type": "disk",
    "shard_num": 1,
    "mappings": {
        "properties": {
            "title": {
                "type": "text",
                "index": true,
                "store": true,
                "highlightable": true,
                "analyzer": "gse_search",
                "search_analyzer": "gse_standard"
            },
            "ownerName": {
                "type": "text",
                "index": true,
                "store": true,
                "highlightable": true,
                "analyzer": "gse_search",
                "search_analyzer": "gse_standard"
            },
            "author": {
                "type": "text",
                "index": true,
                "store": true,
                "highlightable": true,
                "analyzer": "gse_search",
                "search_analyzer": "gse_standard"
            },
            "channelDesc": {
                "type": "text",
                "index": true,
                "store": true,
                "highlightable": true,
                "analyzer": "gse_search",
                "search_analyzer": "gse_standard"
            },
            "textChannelDesc": {
                "type": "text",
                "index": true,
                "store": true,
                "highlightable": true,
                "analyzer": "gse_search",
                "search_analyzer": "gse_standard"
            },
            "id": {
                "type": "text",
                "store": true,
                "index": false
            },
            "link": {
                "type": "text",
                "store": true,
                "index": false
            },
            "feedLink": {
                "type": "text",
                "store": true,
                "index": false
            },
            "language": {
                "type": "keyword",
                "store": true,
                "index": false
            },
            "copyright": {
                "type": "text",
                "store": true,
                "index": false
            },
            "imageUrl": {
                "type": "text",
                "store": true,
                "index": false
            },
            "ownerEmail": {
                "type": "text",
                "store": true,
                "index": false
            },
            "feedType": {
                "type": "text",
                "store": true,
                "index": false
            },
            "categories": {
                "type": "text",
                "store": true,
                "index": true
            }
        }
    }
}

`
