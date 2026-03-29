import { createVuetify } from 'vuetify'
import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

// 浅蓝色主题配置
export default createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'statboxLight',
    themes: {
      statboxLight: {
        dark: false,
        colors: {
          primary: '#5B9BD5',
          'primary-light': '#8DB8E8',
          'primary-dark': '#3A7BC8',
          secondary: '#7A8B99',
          accent: '#F4A261',
          surface: '#FFFFFF',
          'surface-variant': '#F5F7FA',
          background: '#F8FAFB',
          'background-secondary': '#EEF2F6',
          'on-surface': '#2C3E50',
          'on-background': '#2C3E50',
          error: '#EF5350',
          info: '#42A5F5',
          success: '#4CAF50',
          warning: '#FF9800',
          border: '#E1E8ED',
          'border-light': '#EEF2F6',
          divider: '#D8E0E7',
        }
      }
    }
  },
  defaults: {
    VBtn: {
      style: 'text-transform: none;',
      rounded: 'md',
      fontWeight: 'medium'
    },
    VCard: {
      rounded: 'lg',
      elevation: 1
    },
    VTextField: {
      rounded: 'md',
      variant: 'outlined',
      density: 'compact'
    }
  }
})
