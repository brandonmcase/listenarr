import { useQuery } from 'react-query'
import { Box, Grid, Paper, Typography, CircularProgress } from '@mui/material'
import { libraryService } from '@/services/library'
import { downloadService } from '@/services/download'

export function Dashboard() {
  const { data: library, isLoading: libraryLoading } = useQuery(
    'library',
    libraryService.getAll
  )

  const { data: downloads, isLoading: downloadsLoading } = useQuery(
    'downloads',
    downloadService.getAll
  )

  if (libraryLoading || downloadsLoading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
      </Box>
    )
  }

  const totalBooks = library?.length || 0
  const activeDownloads = downloads?.filter((d) => d.status === 'downloading').length || 0
  const completedDownloads = downloads?.filter((d) => d.status === 'completed').length || 0

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Dashboard
      </Typography>

      <Grid container spacing={3} sx={{ mt: 2 }}>
        <Grid item xs={12} sm={6} md={3}>
          <Paper sx={{ p: 2 }}>
            <Typography color="textSecondary" gutterBottom>
              Total Books
            </Typography>
            <Typography variant="h4">{totalBooks}</Typography>
          </Paper>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Paper sx={{ p: 2 }}>
            <Typography color="textSecondary" gutterBottom>
              Active Downloads
            </Typography>
            <Typography variant="h4">{activeDownloads}</Typography>
          </Paper>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Paper sx={{ p: 2 }}>
            <Typography color="textSecondary" gutterBottom>
              Completed Downloads
            </Typography>
            <Typography variant="h4">{completedDownloads}</Typography>
          </Paper>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Paper sx={{ p: 2 }}>
            <Typography color="textSecondary" gutterBottom>
              Processing Queue
            </Typography>
            <Typography variant="h4">0</Typography>
          </Paper>
        </Grid>
      </Grid>
    </Box>
  )
}

