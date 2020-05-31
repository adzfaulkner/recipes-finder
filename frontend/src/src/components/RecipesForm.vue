<template>
  <div>
  <div v-if="showResults|showLoading|showNoResults">
    <transition name="modal">
      <div class="modal-mask">
        <div class="modal-wrapper">
          <div class="modal-dialog" role="document">
            <div class="modal-content">
                  <div class="modal-header">
                      <h5 class="modal-title" id="exampleModalLongTitle">Search Results</h5>
                      <button type="button" class="close" @click="closeModal" aria-label="Close">
                          <span aria-hidden="true">&times;</span>
                      </button>
                  </div>
                  <div v-if="showLoading" class="modal-body spinner">
                      <div class="text-center">
                          <div class="spinner-border" role="status">
                            <span class="sr-only">Loading...</span>
                          </div>
                        </div>
                  </div>
                  <div v-else-if="showResults" class="modal-body results">
                    <div class="container">
                      <h3 class="mb-4">Good matches</h3>
                      <ul class="list-unstyled good-matches">
                          <li v-for="result in goodResults" v-bind:key="result.title" class="media mb-4">
                              <svg class="bd-placeholder-img mr-3" width="64" height="64" xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="xMidYMid slice" focusable="false" role="img" aria-label="Placeholder: 64x64">
                                  <rect width="100%" height="100%" fill="#868e96"></rect>
                                  <text x="12%" y="50%" fill="#dee2e6" dy=".3em">64x64</text>
                              </svg>
                              <div class="media-body">
                                <h5 class="mt-0 mb-1 title">{{result.title}}</h5>
                                <p>Serves {{result.serves}} <span v-if="result.serves == 1">person</span><span v-else>people</span></p>
                                <ul v-for="ingredient in result.ingredients" v-bind:key="ingredient.item" class="ingredients">
                                  <li>{{ingredient.quantity}}{{ingredient.measure}} {{ingredient.item}}</li>
                                </ul>
                              </div>
                            </li>
                      </ul>
                    </div>
                    <div v-if="partialResults.length > 0" class="container">
                      <h3 class="mb-4">Partial matches</h3>
                      <ul class="list-unstyled partial-matches">
                          <li v-for="result in partialResults" v-bind:key="result.title" class="media mb-4">
                              <svg class="bd-placeholder-img mr-3" width="64" height="64" xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="xMidYMid slice" focusable="false" role="img" aria-label="Placeholder: 64x64">
                                  <rect width="100%" height="100%" fill="#868e96"></rect>
                                  <text x="12%" y="50%" fill="#dee2e6" dy=".3em">64x64</text>
                              </svg>
                              <div class="media-body">
                                <h5 class="mt-0 mb-1 title">{{result.title}}</h5>
                                <p>Serves {{result.serves}} <span v-if="result.serves == 1">person</span><span v-else>people</span></p>
                                <ul v-for="ingredient in result.ingredients" v-bind:key="ingredient.item" class="ingredients">
                                  <li>{{ingredient.quantity}}{{ingredient.measure}} {{ingredient.item}}</li>
                                </ul>
                              </div>
                            </li>
                      </ul>
                    </div>  
                  </div>
                  <div v-else-if="showNoResults" class="modal-body zero-results">
                      <ul class="list-unstyled result">
                          <li class="media" style="justify-content: center;">
                              No results founds. Please try again.
                          </li>
                      </ul>
                  </div>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </div>   
      <div class="container">
          <div class="jumbotron mt-5">
              <h1 class="display-4">Welcome!</h1>
              <p class="lead">Enter the ingredients at your disposal and let us find you the recipes!</p>
          </div>
          <form ref="form" @submit.prevent="processForm" action="#" method="post">
              <div class="inputs">
                  <div v-for="(ingredient, index) in ingredients" v-bind:key="index" class="form-row mb-2">
                      <div class="col-9">
                          <input type="text" placeholder="Ingredient" class="form-control ingredient-input" v-model="ingredient.name">
                      </div>
                      <div class="col-2">
                          <input type="number" placeholder="Quantity" class="form-control quantity-input" v-model="ingredient.quantity" number>
                      </div>
                      <div class="col-1">
                          <button type="button" class="btn btn-danger float-right btn-delete" @click="deleteIngredient(index)">Delete</button>
                      </div> 
                  </div>   
              </div>
              <div class="buttons">
                  <div class="form-row">
                      <div class="col-6">
                          <button type="button" class="btn btn-primary btn-add" @click="addIngredient">Add</button>
                      </div>
                      <div class="col-6">
                          <button type="submit" class="btn btn-primary btn-submit float-right">Submit</button>
                      </div>  
                  </div>     
              </div>
          </form>
      </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'RecipesForm',
  data: () => {
    return { 
      ingredients: [{ name: '', quantity: '' }], 
      showResults: false,
      showLoading: false,
      showNoResults: false,
      goodResults: [],
      partialResults: [],
    }
  },
  methods: {
    closeModal() {
      this.showResults = false;
      this.showLoading = false;
      this.showNoResults = false;
    },
    addIngredient() {
      this.ingredients = [...this.ingredients, { name: '', quantity: '' }]
    },
    deleteIngredient(i) {
      let ingredients = this.ingredients;

      if (ingredients.length < 2) {
        this.ingredients = [{ name: '', quantity: '' }];
        return;
      }

      ingredients = ingredients.slice(0, i)
        .concat(ingredients.slice(i + 1, ingredients.length));

      this.ingredients = ingredients;
    },
    processForm() {
      const formValues = this.ingredients.map(ingredient => { return { ingredient: ingredient.name, quantity: parseInt(ingredient.quantity, 10)}});

      this.showLoading = true;

      let apiUrl;

      const {VUE_APP_API_URL, VUE_APP_SERVER_PORT} = process.env;

      //indicates that this is being executed as part of a CI test
      if (VUE_APP_API_URL !== "") {
        apiUrl = VUE_APP_API_URL;
      } else {
        const { location: { hostname } } = window;
        apiUrl = `${hostname}:${VUE_APP_SERVER_PORT}`;
      }

      axios
        .post(`http://${apiUrl}/search`, formValues)
        .then(({data: {id}}) => {
          pollForResults(() => axios.get(`http://${apiUrl}/results?id=${id}`, formValues), 1000000, 2000)
            .then(({good_results, partial_results}) => {
              this.showLoading = false;

              if (good_results.length < 1 && partial_results.length < 1) {
                this.showNoResults = true;
                return;
              }

              this.goodResults = good_results;
              this.partialResults = partial_results;
              this.showResults = true;
            });
        })
        .catch(error => {
          console.log(error.response)
          this.showLoading = false;
          this.showNoResults = true;
        });
    }
  }
}

function pollForResults(fn, timeout, interval) {
    var endTime = Number(new Date()) + (timeout || 2000);
    interval = interval || 100;

    var checkCondition = function(resolve, reject) { 
        var ajax = fn();
        // dive into the ajax promise
        ajax.then( function(response){
            // If the condition is met, we're done!
            if(response.data.success === true) {
                resolve(response.data);
            }

            // If the condition isn't met but the timeout hasn't elapsed, go again
            else if (Number(new Date()) < endTime) {
                setTimeout(checkCondition, interval, resolve, reject);
            }
            // Didn't match and too much time, reject!
            else {
                reject(new Error('timed out for ' + fn + ': ' + arguments));
            }
        });
    };

    return new Promise(checkCondition);
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
form .buttons {
    margin: 30px 0;
}

.modal-mask {
  position: fixed;
  z-index: 9998;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, .5);
  display: table;
  transition: opacity .3s ease;
  overflow-y: auto;
}

.modal-wrapper {
  display: table-cell;
  vertical-align: middle;
}

/* Important part */
.modal-dialog{
    overflow-y: initial !important
}
.modal-body{
    height: 500px;
    overflow-y: auto;
}
</style>
