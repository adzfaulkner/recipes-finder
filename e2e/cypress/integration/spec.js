// enables intelligent code completion for Cypress commands
// https://on.cypress.io/intelligent-code-completion
/// <reference types="Cypress" />

context('Recipe finder test', () => {
  beforeEach(() => {
    // https://on.cypress.io/visit
    cy.visit('/')
  })

  it('delete the only row', function () {  
    cy.get('.ingredient-input')
      .type('tomato')
    cy.get('.quantity-input')
      .type('450')        
    cy.get('.form-row .btn-delete')
      .click()      
    cy.get('.inputs .form-row')
      .should('have.length', 1) 
    cy.get('.ingredient-input')
      .should('have.value', "")   
      cy.get('.quantity-input')
      .should('have.value', "")                               
  })

  it('add 3 then delete a row', function () {
    cy.get('.btn-add')
      .click()
      .click()   
    cy.get('.form-row:nth-child(2) .btn-delete')
      .click()   
    cy.get('.inputs .form-row')
      .should('have.length', 2)                     
  })

  it('search that results in no results', function () {
    cy.get('.ingredient-input')
      .type('invalid')
    cy.get('.quantity-input')
      .type('100')      
    cy.get('.btn-submit')
      .click()
    cy.get('.zero-results')
      .should('be.visible')             
  })

  it('search with 1 criteria', function () {
    cy.get('.ingredient-input')
      .type('tomato')
    cy.get('.quantity-input')
      .type('450')      
    cy.get('.btn-submit')
      .click()
    cy.get('.spinner')
      .should('be.visible')      
    cy.get('.results')
      .should('be.visible')
    cy.get('.good-matches li.media')
      .should('have.length', 6)    
    cy.get('.partial-matches li.media')
      .should('have.length', 0)         
  })

  it('search with 2 criteria', function () {
    cy.get('.ingredient-input')
      .type('tomato')
    cy.get('.quantity-input')
      .type('450')      
    cy.get('.btn-add')
      .click()
    cy.get('.form-row:nth-child(2) .ingredient-input')
      .type('spinach')
    cy.get('.form-row:nth-child(2) .quantity-input')
      .type('100') 
      cy.get('.btn-submit')
      .click()
    cy.get('.spinner')
      .should('be.visible')      
    cy.get('.results')
      .should('be.visible')
    cy.get('.good-matches li.media')
      .should('have.length', 2)    
    cy.get('.partial-matches li.media')
      .should('have.length', 7)   
  })  
})
