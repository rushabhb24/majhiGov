const CACHE_NAME = 'majhigov-pwa-v1';
const ASSETS_TO_CACHE = [
  '/',
  '/index.html',
  '/favicon.svg',
  '/icons.svg'
];

// Install Event
self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => {
      return cache.addAll(ASSETS_TO_CACHE);
    })
  );
  self.skipWaiting();
});

// Activate Event
self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((keys) => {
      return Promise.all(
        keys.map((key) => {
          if (key !== CACHE_NAME) {
            return caches.delete(key);
          }
        })
      );
    })
  );
  self.clients.claim();
});

// Fetch Interception
self.addEventListener('fetch', (event) => {
  // Only intercept GET requests and http/https protocols
  if (event.request.method !== 'GET') {
    return;
  }

  const url = new URL(event.request.url);
  if (url.protocol !== 'http:' && url.protocol !== 'https:') {
    return;
  }

  // Handle API Requests separately to avoid static asset cache-first interference
  if (url.pathname.startsWith('/api/')) {
    // Network-first strategy for dynamic dynamic schemes, jobs, and saved items
    if (url.pathname.includes('/api/schemes') || url.pathname.includes('/api/jobs') || url.pathname.includes('/api/user/saved')) {
      event.respondWith(
        fetch(event.request)
          .then((response) => {
            // Only cache successful dynamic GET requests
            if (response.status === 200) {
              const responseClone = response.clone();
              caches.open(CACHE_NAME).then((cache) => {
                cache.put(event.request, responseClone);
              });
            }
            return response;
          })
          .catch(() => {
            // Serve from cache if offline
            return caches.match(event.request);
          })
      );
    } else {
      // Let other API calls (e.g. auth login/register/logout/profile/notifications) pass through to network
      // Do NOT call event.respondWith, allowing browser to handle the request normally
      return;
    }
    return;
  }

  // Cache-first strategy for static assets
  event.respondWith(
    caches.match(event.request).then((cachedResponse) => {
      if (cachedResponse) {
        return cachedResponse;
      }
      return fetch(event.request).then((response) => {
        // Caching newly requested static assets belonging to our application
        if (response.status === 200 && 
            (event.request.destination === 'image' || 
             event.request.destination === 'font' || 
             event.request.destination === 'style' || 
             event.request.destination === 'script')) {
          const responseClone = response.clone();
          caches.open(CACHE_NAME).then((cache) => {
            cache.put(event.request, responseClone);
          });
        }
        return response;
      });
    })
  );
});
