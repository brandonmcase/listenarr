import { Routes, Route } from 'react-router-dom'
import { Box } from '@mui/material'

import { Layout } from './components/Layout/Layout'
import { Dashboard } from './pages/Dashboard/Dashboard'
import { Library } from './pages/Library/Library'
import { Downloads } from './pages/Downloads/Downloads'
import { Processing } from './pages/Processing/Processing'
import { Search } from './pages/Search/Search'
import { Settings } from './pages/Settings/Settings'

function App() {
  return (
    <Layout>
      <Box component="main" sx={{ flexGrow: 1, p: 3 }}>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/library" element={<Library />} />
          <Route path="/downloads" element={<Downloads />} />
          <Route path="/processing" element={<Processing />} />
          <Route path="/search" element={<Search />} />
          <Route path="/settings" element={<Settings />} />
        </Routes>
      </Box>
    </Layout>
  )
}

export default App

