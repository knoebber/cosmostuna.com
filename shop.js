const getProductURL = "https://p4b43mv7al.execute-api.us-west-2.amazonaws.com/dev/products"

fetch(getProductURL)
  .then((response) => {
    if (!response.ok) {
      throw 'failed to get products';
    }
    return response.json();
  })
  .then((json) => {
    console.log(json);
  })
  .catch((err) => console.log(err));
