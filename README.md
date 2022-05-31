# About Project
基于 Postgresql 和 Postgis 的矢量瓦片免切片服务器，支持多坐标数据源，支持切片缓存。

# How To Run
```go
go run main.go
```

# How To Use
1. 配置数据库信息，启动后台服务，矢量瓦片服务为：`http://${host}:${port}/service/${tablename}/{z}/{x}/{y}/pbf`，瓦片中图层名为${tablename}
2. OpenLayers
```js
import MVT from 'ol/format/MVT.js';
import VectorTileLayer from 'ol/layer/VectorTile.js';
import VectorTileSource from 'ol/source/VectorTile.js';
import {Stroke, Style} from 'ol/style.js';
const vtLayer = new VectorTileLayer({
  declutter: false,
  source: new VectorTileSource({
    format: new MVT(),
    url: 'http://${host}:${port}/service/${tablename}/{z}/{x}/{y}/pbf'  // ${host} ${port} ${tablename}需要替换
  }),
  style: new Style({
      stroke: new Stroke({
        color: 'red',
        width: 1
      })
  })
});
```
3. Mapbox
```js
const map = new mapboxgl.Map({
    'container': 'map',
    'zoom': 14,
    'center': [113, 23], // Guangzhou
    'style': {
        'version': 8,
        'sources': {
            'postgis-tiles': {
                'type': 'vector',
                'tiles': [
                    "http://${host}:${port}/service/${tablename}/{z}/{x}/{y}/pbf"
                ]
            }

        },
        'layers': [{
            'id': 'postgis-tiles-layer',
            'type': 'line',
            'source': 'postgis-tiles',
            'source-layer': '${tablename}', 
            'minzoom': 0,
            'maxzoom': 22,
            'paint': {
                'line-opacity': 0.7,
                'line-color': 'red',
                'line-width': 1
            }
        }]
    }
});
```