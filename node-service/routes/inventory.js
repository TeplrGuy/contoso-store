var express = require('express');
var router = express.Router();
const axios = require('axios').default;
const inventoryService = process.env.INVENTORY_SERVICE_NAME || 'go-app';
const daprPort = process.env.DAPR_HTTP_PORT || 3500;

//use dapr http proxy (header) to call inventory service with normal /inventory route URL in axios.get call
const daprSidecar = `http://localhost:${daprPort}`
//const daprSidecar = `http://localhost:${daprPort}/v1.0/invoke/${inventoryService}/method`

/* GET inventory item by ID */
router.get('/', async function(req, res, next) {
  try {
    var data = await axios.get(`${daprSidecar}/inventory?id=${req.query.id}`, {
      headers: {'dapr-app-id' : `${inventoryService}`} //sets app name for service discovery
    });

    // Return JSON response directly
    res.json(data.data);
  } catch (error) {
    res.status(error.response?.status || 500).json({ 
      error: error.response?.data || 'Failed to retrieve inventory' 
    });
  }
});

/* POST create inventory item */
router.post('/', async function(req, res, next) {
  try {
    var data = await axios.post(`${daprSidecar}/inventory`, req.body, {
      headers: {
        'dapr-app-id': `${inventoryService}`,
        'Content-Type': 'application/json'
      }
    });

    res.status(201).json(data.data);
  } catch (error) {
    res.status(error.response?.status || 500).json({ 
      error: error.response?.data || 'Failed to create inventory' 
    });
  }
});

module.exports = router;
