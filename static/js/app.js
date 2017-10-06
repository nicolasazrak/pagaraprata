var isProd = location.host !== 'localhost:3000';
var apiBase = isProd ? 'https://storage.bad.mn/pagaraprata' : '';
var id = isProd ? 'H0gYbwwzSS6DhX4fCn4oUQ' : 'stHJDE6zSPOsKb5SbSZEZw';
var secret = null;

(function() {

  /* Routing code */

  var splitted = window
    .location
    .hash
    .replace("#", "")
    .replace(/^\//, "")
    .split("/")
    ;

  if (splitted.length === 1 && splitted[0].length > 5) {
    id = splitted[0];
  }

  if (splitted.length === 3 && splitted[1] === "edit") {
    id = splitted[0];
    secret = splitted[2];
  }
})()

new Vue({
  el: '#app',
  data: {
    costs: [],
    debtors: []
  },
  mounted: function () {
    fetch(apiBase + '/api/debts/' + id)
      .then(function (response) { return response.json(); })
      .then(function (data) {
        this.costs = data.costs;
        this.debtors = data.debtors;
      }.bind(this))
      ;
  },
  methods: {
    getDate() {
      var today = new Date();
      var dd = today.getDate();
      var mm = today.getMonth() + 1; //January is 0!
      var yyyy = today.getFullYear();
      if (dd < 10) {
        dd = '0' + dd;
      }
      if (mm < 10) {
        mm = '0' + mm;
      }
      return dd + '/' + mm + '/' + yyyy;
    },
    removeCost(cost) {
      if (confirm('Seguro?')) {
        this.costs = _.without(this.costs, cost);
        this.save();
      }
    },
    addCost() {
      var description = prompt('Que es?');
      var quantity = prompt('¿Cuantas cuotas?');
      var value = prompt('¿Cuanto cada cuota?');

      if (!value || isNaN(value)) { return; }

      this.costs.push({
        value: value,
        quantity: quantity,
        description: description
      });
      this.save();
    },
    removePaymentFor(name, payment) {
      if (confirm('Seguro?')) {
        var debtor = this.findDebtor(name);
        debtor.payments = _.without(debtor.payments, payment);
        this.save();
      }
    },
    addPaymentFor(name) {
      var value = prompt('¿Cuanto?');
      if (!value || isNaN(value)) { return; }

      this.findDebtor(name).payments.push({
        quantity: parseInt(value),
        date: this.getDate(),
      });
      this.save();
    },
    total: function () {
      return _.sum(_.map(this.costs, (function (cost) {
        return cost.value * cost.quantity;
      })));
    },
    findDebtor: function (name) {
      return _.find(this.debtors, function (d) {
        return d.name === name;
      }.bind(this));
    },
    paidBy: function (name) {
      return _.sum(_.map(this.findDebtor(name).payments, function (p) { return p.quantity }.bind(this)));
    },
    save: function () {
      fetch(apiBase + '/api/debts/' + id + '/' + secret, {
        method: 'put',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ costs: this.costs, debtors: this.debtors })
      });
    }
  },
  computed: {
    remaining: function () {
      return _.reduce(this.debtors, function (acumm, currentDebtor) {
        acumm[currentDebtor.name] = this.total() - this.paidBy(currentDebtor.name);
        return acumm;
      }.bind(this), {});
    },
    canEdit: function () {
      return secret !== null;
    }
  }
})