import { useQuery } from 'react-query'
import { Box, Typography, CircularProgress } from '@mui/material'
import { libraryService } from '@/services/library'

export function Library() {
  const { data: library, isLoading } = useQuery('library', libraryService.getAll)

  if (isLoading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
      </Box>
    )
  }

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Library
      </Typography>
      <Typography variant="body1" color="textSecondary">
        Library management - to be implemented
      </Typography>
      <Typography variant="body2" sx={{ mt: 2 }}>
        Total items: {library?.length || 0}
      </Typography>
    </Box>
  )
}

