import 'ol/ol.css';
import TileLayer from 'ol/layer/Tile';
import OSM from 'ol/source/OSM';
import Map from 'ol/Map.js';
import View from 'ol/View.js';
import MVT from 'ol/format/MVT.js';
import VectorTileLayer from 'ol/layer/VectorTile.js';
import VectorTileSource from 'ol/source/VectorTile.js';
import {Stroke, Style} from 'ol/style.js';

const vtLayer = new VectorTileLayer({
  declutter: false,
  source: new VectorTileSource({
    format: new MVT(),
    url: 'http://${host}:${port}/service/${tablename}/{z}/{x}/{y}/pbf'
  }),
  style: new Style({
      stroke: new Stroke({
        color: 'red',
        width: 1
      })
  })
});

const tLayer = new TileLayer({
      source: new OSM()
    });


const map = new Map({
  target: 'map',
  layers: [
    tLayer,
    vtLayer
  ],
  view: new View({
    center: [0, 0],
    zoom: 14
  })
});

