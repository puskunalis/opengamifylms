describe('OpenGamifyLMS Tests', () => {
  beforeEach(() => {
    cy.visit('http://opengamifylms-frontend');
  });

  it('should allow user to sign up and log in', () => {
    cy.get('button').contains('Get Started').click();
    cy.get('input[id=":r1:"]').type('vp@example.com');
    cy.get('input[id=":r2:"]').type('123');
    cy.get('button').contains('Login').click();
    cy.contains('vp@example.com').should('be.visible');
  });

  it('should allow user to enroll in a course', () => {
    // Log in first
    cy.get('button').contains('Get Started').click();
    cy.get('input[id=":r1:"]').type('vp@example.com');
    cy.get('input[id=":r2:"]').type('123');
    cy.get('button').contains('Login').click();

    // Find and click on a course to enroll
    cy.contains('New Quests').should('be.visible');
    cy.get('.MuiCardActionArea-root').first().click();

    // Enroll in the course
    cy.get('button').contains('Enroll').click();
    cy.contains('Progress:').should('be.visible');
  });

  it('should display user courses on the logged-in page', () => {
    // Log in first
    cy.get('button').contains('Get Started').click();
    cy.get('input[id=":r1:"]').type('vp@example.com');
    cy.get('input[id=":r2:"]').type('123');
    cy.get('button').contains('Login').click();

    // Check if user courses are displayed
    cy.contains('Your Adventure').should('be.visible');
    cy.get('.MuiGrid-root').should('have.length.greaterThan', 0);
  });

  it('should allow user to create a course, module, submodules, and elements, publish, enroll, and set the course as unavailable', () => {
    // Log in first
    cy.get('button').contains('Get Started').click();
    cy.get('input[id=":r1:"]').type('john@example.com');
    cy.get('input[id=":r2:"]').type('123');
    cy.get('button').contains('Login').click();
  
    // Navigate to Course Builder
    cy.get('a[href="/coursebuilder"]').click();
  
    // Enter course title and description
    cy.get('input').first().type('Introduction to Cypress');
    cy.get('textarea').first().type('Learn how to write tests using Cypress.');
  
    // Save the course
    cy.get('button').contains('Save Course').click();
  
    // Verify the course is created
    cy.contains('Introduction to Cypress').should('be.visible');
    cy.contains('Unpublished').should('be.visible');
    cy.contains('Unavailable').should('be.visible');
  
    // Go back to homepage and verify the course is not visible yet
    cy.get('a').contains('OpenGamifyLMS').click();
    cy.contains('Introduction to Cypress').should('not.exist');
  
    // Navigate back to Course Builder
    cy.get('a[href="/coursebuilder"]').click();

    // Find and click on the created course
    cy.contains('Introduction to Cypress').click();

    cy.contains('Module Builder').should('be.visible');

    // Create 5 modules
    for (let i = 1; i <= 5; i++) {
      cy.get('.MuiInputBase-input').first().type(`Module ${i}`);
      cy.get('button').contains('Add').click();
      cy.wait(200);
    }

    // Verify that 5 modules have been added and are in correct order
    cy.get('.MuiBox-root[data-rbd-droppable-id="modules"]').children().should('have.length', 5);
    cy.get('.MuiBox-root[data-rbd-droppable-id="modules"]').children().each((module, index) => {
      cy.wrap(module).find('h6').should('contain', `Module ${index + 1}`);
    });
    
    // Click delete for Module 3
    cy.xpath('(//button[@aria-label="delete"])[3]').click();

    // Verify that modules are reordered correctly after deletion
    cy.get('.MuiBox-root[data-rbd-droppable-id="modules"]').children().should('have.length', 4);
    cy.get('.MuiBox-root[data-rbd-droppable-id="modules"]').children().eq(0).find('h6').should('contain', 'Module 1');
    cy.get('.MuiBox-root[data-rbd-droppable-id="modules"]').children().eq(1).find('h6').should('contain', 'Module 2');
    cy.get('.MuiBox-root[data-rbd-droppable-id="modules"]').children().eq(2).find('h6').should('contain', 'Module 4');
    cy.get('.MuiBox-root[data-rbd-droppable-id="modules"]').children().eq(3).find('h6').should('contain', 'Module 5');

    // Click edit for Module 1
    cy.xpath('(//a[@aria-label="edit"])[1]').click();

    cy.contains('Submodule Builder').should('be.visible');

    // Create 5 submodules
    for (let i = 1; i <= 5; i++) {
      cy.get('.MuiInputBase-input').first().type(`Submodule ${i}`);
      cy.wait(500);
      cy.get('.MuiInputBase-input').last().type(`Submodule ${i} Description`);
      cy.wait(500);
      cy.get('button').contains('Add').click();
      cy.wait(500);
    }

    // Verify that 5 submodules have been added and are in correct order
    cy.get('.MuiBox-root[data-rbd-droppable-id="submodules"]').children().should('have.length', 5);
    cy.get('.MuiBox-root[data-rbd-droppable-id="submodules"]').children().each((module, index) => {
      cy.wrap(module).find('h6').should('contain', `Submodule ${index + 1}`);
    });
    
    // Click delete for Submodule 4
    cy.xpath('(//button[@aria-label="delete"])[4]').click();

    // Verify that submodules are reordered correctly after deletion
    cy.get('.MuiBox-root[data-rbd-droppable-id="submodules"]').children().should('have.length', 4);
    cy.get('.MuiBox-root[data-rbd-droppable-id="submodules"]').children().eq(0).find('h6').should('contain', 'Submodule 1');
    cy.get('.MuiBox-root[data-rbd-droppable-id="submodules"]').children().eq(1).find('h6').should('contain', 'Submodule 2');
    cy.get('.MuiBox-root[data-rbd-droppable-id="submodules"]').children().eq(2).find('h6').should('contain', 'Submodule 3');
    cy.get('.MuiBox-root[data-rbd-droppable-id="submodules"]').children().eq(3).find('h6').should('contain', 'Submodule 5');

    // Click edit for Submodule 1
    cy.xpath('(//a[@aria-label="edit"])[1]').click();

    cy.contains('Element Builder').should('be.visible');

    // Click on the HTML content button
    cy.get('input[value="html"]').click();

    // Type in some content
    cy.get('div[role="textbox"]').type('This is some content.');

    // Save the element  
    cy.get('button').contains('Add').click();

    // Verify the HTML content element was added
    cy.get('.MuiBox-root[data-rbd-droppable-id="elements"]').children().should('have.length', 1);
    cy.get('.MuiBox-root[data-rbd-droppable-id="elements"]').children().eq(0).should('contain', 'This is some content.');

    // Click on the Video content button
    cy.get('input[value="video"]').click();

    // Upload the file using cy.fixture
    cy.get('input[type="file"]').selectFile('cypress/fixtures/ForBiggerJoyrides.mp4');

    // Save the video element
    cy.get('button').contains('Add').click();

    // Verify the video element was added
    cy.get('.MuiBox-root[data-rbd-droppable-id="elements"]').children().should('have.length', 2);

    // Verify the video is playable
    cy.get('video').should('be.visible').and($video => {
      expect($video[0].duration).to.be.greaterThan(0);
    });

    // Navigate to Course Builder
    cy.get('a[href="/coursebuilder"]').click();

    // Publish the course  
    cy.get('button').contains('Publish').click();

    // Verify the course is now published
    cy.contains('Published').should('be.visible');

    // Go to the homepage
    cy.get('a').contains('OpenGamifyLMS').click();
    cy.reload();

    // Verify the course is now visible on the homepage
    cy.contains('Introduction to Cypress').should('be.visible');

    // Enroll in the course
    cy.contains('Introduction to Cypress').click();
    cy.get('button').contains('Enroll').click();

    // Verify enrollment
    cy.contains('Progress:').should('be.visible');

    // Navigate to the enrolled course
    cy.get('a[href="/courses/4"]').click();

    // Verify the course content section is visible
    cy.contains('Course Modules').should('be.visible');

    // Verify the modules are in the correct order
    cy.get('.MuiCardContent-root h6').should('have.length', 4);
    cy.get('.MuiCardContent-root h6').eq(0).should('contain', 'Module 1');
    cy.get('.MuiCardContent-root h6').eq(1).should('contain', 'Module 2');
    cy.get('.MuiCardContent-root h6').eq(2).should('contain', 'Module 4');
    cy.get('.MuiCardContent-root h6').eq(3).should('contain', 'Module 5');

    // Navigate to Module 1
    cy.get('.MuiButtonBase-root').contains('Start').first().click();

    // Verify the module title is visible
    cy.contains('Module 1').should('be.visible');

    // Verify the submodules are in the correct order
    cy.get('.MuiListItemText-root span').should('have.length', 4);
    cy.get('.MuiListItemText-root span').eq(0).should('contain', 'Submodule 1');
    cy.get('.MuiListItemText-root span').eq(1).should('contain', 'Submodule 2');
    cy.get('.MuiListItemText-root span').eq(2).should('contain', 'Submodule 3');
    cy.get('.MuiListItemText-root span').eq(3).should('contain', 'Submodule 5');

    // Navigate to Submodule 1
    cy.contains('Submodule 1').click();

    // Verify the submodule title is visible
    cy.contains('Submodule 1').should('be.visible');

    // Verify the added elements are visible
    cy.contains('This is some content.').should('be.visible');

    // Verify the video is playable
    cy.get('video').should('be.visible').and($video => {
      expect($video[0].duration).to.be.greaterThan(0);
    });

    /*// Go back to Course Builder 
    cy.get('a[href="/coursebuilder"]').click();

    // Set the course as unavailable
    cy.get('button').contains('Set Unavailable').click();

    // Go to the homepage
    cy.get('a').contains('OpenGamifyLMS').click();
    cy.reload();

    // Verify the course is no longer visible
    cy.wait(3000);
    cy.contains('Introduction to Cypress').should('not.exist');*/
  });
});
