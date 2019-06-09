;(function() {
  const WhmIndex = window.WhmIndex = {}

  const searchInput = document.querySelector('.search-container input[type="search"]')
  if (searchInput) {
    const handleChange = _.debounce((event) => {
      const value = searchInput.value.trim()
      fetch(`/api/search.json?query=${encodeURIComponent(value)}`)
        .then(response => response.json())
        .then(result => {
          let searchResults = document.querySelector('.search-container .results')
          if (!searchResults) {
            searchResults = document.createElement('div')
            searchResults.className = 'results'
            searchInput.parentElement.appendChild(searchResults)
          }
          if (result.matches.length > 0) {
            const elements = result.matches.map(match => (
              `<li>${match.episode.title}</li>`
            ))
            searchResults.innerHTML = `<ul>${elements.join('')}</ul>`
          } else {
            searchResults.innerHTML = 'No results found.'
          }
        })
    }, 500)
    searchInput.addEventListener('change', handleChange)
    searchInput.addEventListener('input', handleChange)
  }
})();
