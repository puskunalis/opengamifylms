import React from 'react';
import { Container, Typography } from '@mui/material';

function Section({ title, description }) {
 return (
    <section className="mt-8">
      <Container maxWidth="sm">
        <Typography variant="h2" component="h2" gutterBottom>
          {title}
        </Typography>
        <Typography variant="body1" component="p" gutterBottom>
          {description}
        </Typography>
      </Container>
    </section>
 );
}

export default Section;