import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  resolve: {
    alias: {
     crypto: 'crypto-browserify',
     stream: 'stream-browserify'
    }
  },
  optimizeDeps: {
   include: ["crypto-browserify"]
  },
  plugins: [react()],
});
