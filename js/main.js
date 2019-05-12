;(function() {
  const WhmIndex = window.WhmIndex = {}

  const searchInput = document.querySelector('.search-container input[type="search"]')
  if (searchInput) {
    const handleChange = _.debounce((event) => {
      const value = searchInput.value.trim()
      fetch(`/api/search.json?query=${encodeURIComponent(value)}`)
    }, 500)
    searchInput.addEventListener('change', handleChange)
    searchInput.addEventListener('input', handleChange)
  }
})();
