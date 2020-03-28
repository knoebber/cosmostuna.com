<section>
  <?php if(!$prod) :?>
    <div style="text-align: center;"><strong>[TEST MODE]</strong></div>
  <?php endif; ?>
  <div class="image-wrapper">
    <img src="images/new-label.jpg" alt="Cosmo's Tuna Logo">
  </div>
  <br>
  <noscript><h4 style="color:red">Enable scripts to use the shop</h4></noscript>
  <p>
    If you're local, contact Cosmo for a discount.
  </p>
  <hr>
  <div id="product-grid" class="grid-section">
    <em>loading products...</em>
    <div class="spinner"></div>
  </div>
  <form id="shop-form">
    <div id="order-grid" class="grid-section">
      <label for="name">Full Name</label><input id="name" type="text" required>
      <label for="address">Address</label><input id="address" type="text" required>
      <label for="city">City</label><input id="city" type="text" required>
      <label for="state">State</label><input id="state" type="text" required>
      <label for="zip">Zip</label><input id="zip" type="text" pattern="^[0-9]{5}(-[0-9]{4})?$" required>
      <label for="email">Email</label><input id="email" type="email" required>
    </div>
    <div id="submit-row" class="info-row">
      <button type="submit">Next</button>
      <span id="form-errors" for="place-order" role="alert"></label>
    </div>
  </form>
</section>
