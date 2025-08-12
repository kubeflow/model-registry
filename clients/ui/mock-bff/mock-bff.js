const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');

const app = express();
const PORT = 4000;

// Add logging middleware
app.use((req, res, next) => {
  console.log(`${req.method} ${req.url}`);
  next();
});

// Mock UI-specific endpoints that the BFF would normally handle
// Handle both /api and /model-registry/api patterns to be safe
app.get(['/api/v1/user', '/model-registry/api/v1/user'], (req, res) => {
  res.json({
    data: {
      userId: 'user@example.com',
      clusterAdmin: true
    }
  });
});

app.get(['/api/v1/namespaces', '/model-registry/api/v1/namespaces'], (req, res) => {
  res.json({
    data: [
      { name: 'default' },
      { name: 'kubeflow' },
      { name: 'model-registry' }
    ]
  });
});

app.get(['/api/v1/model_registry', '/model-registry/api/v1/model_registry'], (req, res) => {
  console.log('Model registry list request with query:', req.query);
  const response = {
    metadata: {},
    data: [
      {
        name: 'model-registry',
        displayName: 'Model Registry',
        description: 'Default model registry instance',
        serverAddress: 'http://model-registry:8080'
      }
    ]
  };
  console.log('Sending model registry list response:', JSON.stringify(response));
  res.json(response);
});

app.get(['/api/v1/settings/role_bindings', '/model-registry/api/v1/settings/role_bindings'], (req, res) => {
  console.log('Role bindings request with query:', req.query);
  const response = {
    metadata: {},
    data: {
      items: []
    }
  };
  console.log('Sending role bindings response:', JSON.stringify(response));
  res.json(response);
});

app.get(['/api/v1/settings/model_registry', '/model-registry/api/v1/settings/model_registry'], (req, res) => {
  console.log('Model registry settings request with query:', req.query);
  const response = {
    metadata: {},
    data: []
  };
  console.log('Sending response:', JSON.stringify(response));
  res.json(response);
});

app.get(['/api/v1/settings/groups', '/model-registry/api/v1/settings/groups'], (req, res) => {
  res.json({
    data: []
  });
});

app.get(['/api/v1/settings/certificates', '/model-registry/api/v1/settings/certificates'], (req, res) => {
  res.json({
    data: {
      secrets: [],
      configMaps: []
    }
  });
});

// Proxy model-registry API calls to the real model-registry service
// Handle BFF-style endpoints that include the registry name
app.use('/model-registry/api/v1/model_registry/:registryName', createProxyMiddleware({
  target: 'http://model-registry:8080',
  changeOrigin: true,
  pathRewrite: {
    '^/model-registry/api/v1/model_registry/[^/]+': '/api/model_registry/v1alpha3'
  },
  headers: {
    'kubeflow-userid': 'user@example.com'
  },
  logLevel: 'debug',
  selfHandleResponse: true,
  onProxyRes: function(proxyRes, req, res) {
    // Transform model-registry response format to BFF format
    let body = [];
    proxyRes.on('data', function(chunk) {
      body.push(chunk);
    });
    proxyRes.on('end', function() {
      const bodyString = Buffer.concat(body).toString();
      try {
        const originalResponse = JSON.parse(bodyString);
        // Transform from {items: [...], nextPageToken: "", ...} to {metadata: {}, data: {items: [...], ...}}
        const transformedResponse = {
          metadata: {},
          data: originalResponse
        };
        console.log('Transformed response:', JSON.stringify(transformedResponse));
        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify(transformedResponse));
      } catch (error) {
        console.error('Error transforming response:', error);
        res.writeHead(500, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Failed to transform response' }));
      }
    });
  }
}));

app.use('/model-registry/api/model-registry', createProxyMiddleware({
  target: 'http://model-registry:8080',
  changeOrigin: true,
  pathRewrite: {
    '^/model-registry/api/model-registry': '/api/model-registry'
  },
  headers: {
    'kubeflow-userid': 'user@example.com'
  },
  logLevel: 'debug'
}));

app.use('/api/model-registry', createProxyMiddleware({
  target: 'http://model-registry:8080',
  changeOrigin: true,
  headers: {
    'kubeflow-userid': 'user@example.com'
  },
  logLevel: 'debug'
}));

// Catch unmocked API endpoints and log them (must be after proxy middleware)
app.use('/api', (req, res, next) => {
  console.log('UNMOCKED API ENDPOINT:', req.method, req.url);
  res.status(404).json({ error: 'Endpoint not mocked', url: req.url });
});

app.use('/model-registry/api', (req, res, next) => {
  console.log('UNMOCKED MODEL-REGISTRY API ENDPOINT:', req.method, req.url);
  res.status(404).json({ error: 'Endpoint not mocked', url: req.url });
});

app.listen(PORT, '0.0.0.0', () => {
  console.log(`Mock BFF server running on port ${PORT}`);
});
