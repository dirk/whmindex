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
            const elements = result.matches.slice(0, 5).map(match => {
              const href = `/${match.episode.feed}/${match.episode.number}`
              return `<li><a href="${href}">${match.episode.title}</a></li>`
            })
            searchResults.innerHTML = (
              `<ul>${elements.join('')}</ul>` +
              `<div class="all-results"><a href="/search?query=${encodeURIComponent(value)}">` +
                `All results for "${value}" »` +
              '</a></div>'
            )
          } else {
            searchResults.innerHTML = 'No results found.'
          }
        })
    }, 500)
    searchInput.addEventListener('change', handleChange)
    searchInput.addEventListener('input', handleChange)
  }
})();
