// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',
  devtools: { enabled: true },

  modules: [
    '@element-plus/nuxt',
    '@pinia/nuxt',
    '@nuxtjs/tailwindcss',
    '@nuxtjs/i18n',
    '@vite-pwa/nuxt',
  ],

  i18n: {
    locales: [
      { code: 'vi', name: 'Tiếng Việt', file: 'vi.json' },
      { code: 'en', name: 'English', file: 'en.json' },
    ],
    defaultLocale: 'vi',
    langDir: 'locales/',
    strategy: 'prefix_except_default',
  },

  pwa: {
    registerType: 'autoUpdate',
    manifest: {
      name: 'Matcha Comic',
      short_name: 'Matcha',
      theme_color: '#0D0F14',
      icons: [
        { src: 'pwa-192x192.png', sizes: '192x192', type: 'image/png' },
        { src: 'pwa-512x512.png', sizes: '512x512', type: 'image/png' },
      ],
    },
    workbox: {
      navigateFallback: '/',
    },
    devOptions: {
      enabled: true,
      type: 'module',
    },
  },

  elementPlus: {
    importStyle: 'css',
    themes: ['dark'],
  },

  components: [
    {
      path: '~/components',
      pathPrefix: false,
    },
  ],

  css: [
    '~/assets/css/main.css',
    '~/assets/css/animations.css',
    '~/assets/css/3d-card.css',
  ],

  app: {
    head: {
      title: 'Matcha Comic - Đọc Truyện Tranh & Tiểu Thuyết Online',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        {
          name: 'description',
          content: 'Đọc truyện tranh và tiểu thuyết online miễn phí. Cập nhật nhanh nhất, chất lượng cao nhất.',
        },
        { name: 'theme-color', content: '#0D0F14' },
        { property: 'og:type', content: 'website' },
        { property: 'og:title', content: 'Matcha Comic - Đọc Truyện Online' },
      ],
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
        {
          rel: 'preconnect',
          href: 'https://fonts.googleapis.com',
        },
        {
          rel: 'preconnect',
          href: 'https://fonts.gstatic.com',
          crossorigin: '',
        },
        {
          rel: 'stylesheet',
          href: 'https://fonts.googleapis.com/css2?family=Montserrat:wght@400;500;600;700;800;900&family=Inter:wght@300;400;500;600&display=swap',
        },
      ],
    },
  },

  runtimeConfig: {
    public: {
      apiBase: process.env.API_BASE || 'http://localhost:8080/api',
      googleClientId: process.env.NUXT_PUBLIC_GOOGLE_CLIENT_ID || '',
    },
  },

  nitro: {
    compressPublicAssets: true,
  },

  routeRules: {
    '/': { prerender: true },
    '/truyen-tranh/**': { ssr: true },
    '/truyen-chu/**': { ssr: true },
    '/admin/**': { ssr: false },
  },
})
