// Architecture Overview Generator Component
import { useState } from 'react';
import {
  Box,
  Button,
  Card,
  CardContent,
  CircularProgress,
  Grid,
  TextField,
  Typography,
  Snackbar,
  Alert,
} from '@mui/material';
import { API } from '../api';

/**
 * Component for generating architectural overviews of Azure cloud workloads
 */
export default function OverviewGenerator() {
  const [topic, setTopic] = useState('');
  const [overview, setOverview] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [snackbarOpen, setSnackbarOpen] = useState(false);

  /**
   * Handle submission of the form to generate an overview
   */
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!topic.trim()) {
      setError('Please enter a topic to generate an overview');
      setSnackbarOpen(true);
      return;
    }
    
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch(`${API.serverBaseUrl}/api/overview`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ topic }),
      });
      
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Failed to generate overview');
      }
      
      const data = await response.json();
      setOverview(data.overview || 'No overview generated');
    } catch (err: any) {
      setError(err.message || 'Failed to generate overview');
      setSnackbarOpen(true);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box p={2}>
      <Typography variant="h4" gutterBottom>
        Azure Architecture Overview Generator
      </Typography>
      
      <Typography variant="body1" paragraph>
        Enter an Azure workload or solution to generate a comprehensive architectural overview.
      </Typography>
      
      <Card>
        <CardContent>
          <form onSubmit={handleSubmit}>
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  label="Workload or Solution"
                  value={topic}
                  onChange={(e) => setTopic(e.target.value)}
                  placeholder="E.g., Web Application with SQL Database, Microservices Architecture, etc."
                  disabled={loading}
                />
              </Grid>
              <Grid item xs={12}>
                <Button
                  type="submit"
                  variant="contained"
                  color="primary"
                  disabled={loading}
                >
                  {loading ? <CircularProgress size={24} /> : 'Generate Overview'}
                </Button>
              </Grid>
            </Grid>
          </form>
        </CardContent>
      </Card>
      
      {overview && (
        <Card style={{ marginTop: 16 }}>
          <CardContent>              <Typography variant="h5" gutterBottom>
              Azure Architecture Overview: {topic}
            </Typography>
            <Box 
              style={{ 
                whiteSpace: 'pre-wrap',
                backgroundColor: '#f5f5f5',
                padding: 16,
                borderRadius: 4,
              }}
            >
              <Typography variant="body1">{overview}</Typography>
            </Box>
          </CardContent>
        </Card>
      )}
      
      <Snackbar
        open={snackbarOpen}
        autoHideDuration={6000}
        onClose={() => setSnackbarOpen(false)}
      >
        <Alert 
          onClose={() => setSnackbarOpen(false)} 
          severity="error"
        >
          {error}
        </Alert>
      </Snackbar>
    </Box>
  );
}
