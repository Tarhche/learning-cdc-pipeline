package handlers

import (
	"blog-cdc-search/domain"
	"fmt"
)

// generateHomePageHTML generates the main blog page HTML (public view)
func generateHomePageHTML(posts []*domain.Post) string {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Blog - Home</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; }
        .header { text-align: center; margin-bottom: 40px; }
        .header h1 { color: #333; font-size: 2.5em; margin-bottom: 10px; }
        .header p { color: #666; font-size: 1.2em; }
        .admin-link { text-align: center; margin-bottom: 30px; }
        .admin-link a { color: #666; text-decoration: none; font-size: 0.9em; }
        .admin-link a:hover { color: #333; }
        .search-container {
            background: white;
            border-radius: 10px;
            padding: 30px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            margin-bottom: 30px;
        }
        .search-input-container {
            position: relative;
            margin-bottom: 20px;
        }
        .search-input {
            width: 100%;
            padding: 15px 20px;
            font-size: 18px;
            border: 2px solid #ddd;
            border-radius: 25px;
            outline: none;
            transition: border-color 0.3s;
            box-sizing: border-box;
        }
        .search-input:focus {
            border-color: #007bff;
        }
        .search-icon {
            position: absolute;
            right: 20px;
            top: 50%;
            transform: translateY(-50%);
            color: #666;
            font-size: 20px;
        }
        .search-filters {
            display: flex;
            gap: 15px;
            margin-bottom: 20px;
            flex-wrap: wrap;
        }
        .filter-group {
            display: flex;
            flex-direction: column;
            gap: 5px;
        }
        .filter-group label {
            font-weight: bold;
            color: #333;
            font-size: 14px;
        }
        .filter-group select, .filter-group input {
            padding: 8px 12px;
            border: 1px solid #ddd;
            border-radius: 5px;
            font-size: 14px;
        }
        .search-button {
            background-color: #007bff;
            color: white;
            border: none;
            padding: 12px 30px;
            border-radius: 25px;
            font-size: 16px;
            cursor: pointer;
            transition: background-color 0.3s;
        }
        .search-button:hover {
            background-color: #0056b3;
        }
        .search-button:disabled {
            background-color: #ccc;
            cursor: not-allowed;
        }
        .loading {
            text-align: center;
            padding: 20px;
            color: #666;
        }
        .spinner {
            border: 3px solid #f3f3f3;
            border-top: 3px solid #007bff;
            border-radius: 50%;
            width: 30px;
            height: 30px;
            animation: spin 1s linear infinite;
            margin: 0 auto 10px;
        }
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
        .results-info {
            background: white;
            border-radius: 10px;
            padding: 20px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            margin-bottom: 20px;
            text-align: center;
        }
        .results-count {
            font-size: 18px;
            color: #333;
            margin-bottom: 10px;
        }
        .results-query {
            color: #666;
            font-style: italic;
        }
        .no-results {
            text-align: center;
            padding: 40px;
            color: #666;
        }
        .no-results h3 {
            margin-bottom: 10px;
        }
        .posts { display: grid; grid-template-columns: repeat(auto-fill, minmax(350px, 1fr)); gap: 20px; }
        .post { background: white; border-radius: 10px; padding: 20px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); cursor: pointer; transition: transform 0.2s, box-shadow 0.2s; }
        .post:hover { transform: translateY(-2px); box-shadow: 0 4px 20px rgba(0,0,0,0.15); }
        .post img { width: 100%; height: 200px; object-fit: cover; border-radius: 5px; margin-bottom: 15px; }
        .post h2 { color: #333; margin: 0 0 10px 0; font-size: 1.4em; }
        .post .excerpt { color: #666; margin-bottom: 15px; line-height: 1.5; }
        .post .meta { color: #999; font-size: 0.9em; margin-bottom: 15px; }
        .post .read-more { color: #007bff; font-weight: bold; text-align: right; }
        .pagination {
            display: flex;
            justify-content: center;
            gap: 10px;
            margin-top: 30px;
        }
        .pagination button {
            padding: 10px 15px;
            border: 1px solid #ddd;
            background: white;
            color: #333;
            border-radius: 5px;
            cursor: pointer;
            transition: all 0.3s;
        }
        .pagination button:hover:not(:disabled) {
            background: #007bff;
            color: white;
            border-color: #007bff;
        }
        .pagination button:disabled {
            background: #f5f5f5;
            color: #ccc;
            cursor: not-allowed;
        }
        .pagination .current-page {
            background: #007bff;
            color: white;
            border-color: #007bff;
        }
        .view-toggle {
            text-align: center;
            margin-bottom: 20px;
        }
        .view-toggle button {
            padding: 8px 16px;
            margin: 0 5px;
            border: 1px solid #ddd;
            background: white;
            color: #333;
            border-radius: 5px;
            cursor: pointer;
            transition: all 0.3s;
        }
        .view-toggle button.active {
            background: #007bff;
            color: white;
            border-color: #007bff;
        }
        .view-toggle button:hover:not(.active) {
            background: #f8f9fa;
        }
        .post .image-container { position: relative; width: 100%; height: 200px; margin-bottom: 15px; border-radius: 5px; background-color: #f0f0f0; overflow: hidden; }
        .post .image-placeholder { width: 100%; height: 100%; background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%); background-size: 200% 100%; animation: loading 1.5s infinite; }
        .post .lazy-image { width: 100%; height: 100%; object-fit: cover; opacity: 0; transition: opacity 0.3s ease-in-out; }
        .post .lazy-image.loaded { opacity: 1; }
        @keyframes loading {
            0% { background-position: 200% 0; }
            100% { background-position: -200% 0; }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>My Blog</h1>
            <p>Share your thoughts and ideas</p>
        </div>
        
        <div class="admin-link">
            <a href="/dashboard">Admin Dashboard</a>
        </div>
        
        <div class="search-container">
            <div class="search-input-container">
                <input type="text" id="searchInput" class="search-input" placeholder="Search for posts by title, excerpt, or content..." autocomplete="off">
                <div class="search-icon">🔍</div>
            </div>
            
            <div class="search-filters">
                <div class="filter-group">
                    <label for="sortBy">Sort By:</label>
                    <select id="sortBy">
                        <option value="_text_match:desc,created_at:desc">Relevance</option>
                        <option value="created_at:desc">Newest First</option>
                        <option value="created_at:asc">Oldest First</option>
                        <option value="title:asc">Title A-Z</option>
                    </select>
                </div>
                
                <div class="filter-group">
                    <label for="perPage">Results per page:</label>
                    <select id="perPage">
                        <option value="10">10</option>
                        <option value="20">20</option>
                        <option value="50">50</option>
                    </select>
                </div>
            </div>
            
            <button id="searchButton" class="search-button">Search</button>
        </div>
        
        <div class="view-toggle">
            <button id="allPostsBtn" class="active" onclick="showAllPosts()">All Posts</button>
            <button id="searchResultsBtn" onclick="showSearchResults()" style="display: none;">Search Results</button>
        </div>
        
        <div id="loading" class="loading" style="display: none;">
            <div class="spinner"></div>
            <p>Searching...</p>
        </div>
        
        <div id="resultsInfo" class="results-info" style="display: none;">
            <div class="results-count" id="resultsCount"></div>
            <div class="results-query" id="resultsQuery"></div>
        </div>
        
        <div id="noResults" class="no-results" style="display: none;">
            <h3>No results found</h3>
            <p>Try adjusting your search terms or filters</p>
        </div>
        
        <div id="allPosts" class="posts">
`

	if len(posts) == 0 {
		html += `
            <div style="grid-column: 1 / -1; text-align: center; padding: 40px; color: #666;">
                <h3>No posts yet</h3>
                <p>Check back soon for new content!</p>
            </div>
        `
	} else {
		for _, post := range posts {
			html += fmt.Sprintf(`
            <div class="post" onclick="window.location.href='/post/%d'">
                <div class="image-container">
                    <div class="image-placeholder"></div>
                    <img class="lazy-image" data-src="%s" alt="%s" onerror="this.src='https://placehold.co/600x400?text=%s'">
                </div>
                <h2>%s</h2>
                <div class="excerpt">%s</div>
                <div class="meta">Created: %s</div>
                <div class="read-more">Read More →</div>
            </div>
            `, post.ID, post.Image, post.Title, post.Title, post.Title, post.Excerpt, post.CreatedAt.Format("Jan 02, 2006"))
		}
	}

	html += `
        </div>
        
        <div id="searchResults" class="posts" style="display: none;"></div>
        
        <div id="pagination" class="pagination" style="display: none;"></div>
    </div>

    <script>
        let currentPage = 1;
        let currentQuery = '';
        let currentSortBy = '_text_match:desc,created_at:desc';
        let currentPerPage = 10;
        let searchTimeout;
        let isSearchMode = false;

        // Initialize search functionality
        document.addEventListener('DOMContentLoaded', function() {
            const searchInput = document.getElementById('searchInput');
            const searchButton = document.getElementById('searchButton');
            const sortBySelect = document.getElementById('sortBy');
            const perPageSelect = document.getElementById('perPage');

            // Search on input with debouncing
            searchInput.addEventListener('input', function() {
                clearTimeout(searchTimeout);
                searchTimeout = setTimeout(() => {
                    if (this.value.trim().length >= 2) {
                        performSearch();
                    } else if (this.value.trim() === '') {
                        showAllPosts();
                    }
                }, 300);
            });

            // Search on button click
            searchButton.addEventListener('click', performSearch);

            // Search on Enter key
            searchInput.addEventListener('keypress', function(e) {
                if (e.key === 'Enter') {
                    performSearch();
                }
            });

            // Update search when filters change
            sortBySelect.addEventListener('change', function() {
                currentSortBy = this.value;
                if (currentQuery && isSearchMode) {
                    performSearch();
                }
            });

            perPageSelect.addEventListener('change', function() {
                currentPerPage = parseInt(this.value);
                currentPage = 1;
                if (currentQuery && isSearchMode) {
                    performSearch();
                }
            });
        });

        function performSearch() {
            const query = document.getElementById('searchInput').value.trim();
            if (!query) {
                showAllPosts();
                return;
            }

            currentQuery = query;
            isSearchMode = true;
            showLoading();
            showSearchResults();
            
            const searchParams = {
                query: query,
                page: currentPage,
                per_page: currentPerPage,
                sort_by: currentSortBy
            };

            fetch('/api/search', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(searchParams)
            })
            .then(response => response.json())
            .then(data => {
                hideLoading();
                displayResults(data);
            })
            .catch(error => {
                hideLoading();
                console.error('Search error:', error);
                showError('An error occurred while searching. Please try again.');
            });
        }

        function displayResults(data) {
            const resultsContainer = document.getElementById('searchResults');
            const resultsInfo = document.getElementById('resultsInfo');
            const resultsCount = document.getElementById('resultsCount');
            const resultsQuery = document.getElementById('resultsQuery');
            const noResults = document.getElementById('noResults');

            if (!data.results || data.results.length === 0) {
                resultsContainer.innerHTML = '';
                resultsInfo.style.display = 'none';
                noResults.style.display = 'block';
                return;
            }

            // Show results info
            resultsCount.textContent = 'Found ' + data.total + ' result' + (data.total !== 1 ? 's' : '');
            resultsQuery.textContent = 'for "' + data.query + '"';
            resultsInfo.style.display = 'block';
            noResults.style.display = 'none';

            // Display results
            resultsContainer.innerHTML = '';
            data.results.forEach(result => {
                const resultItem = createResultItem(result);
                resultsContainer.appendChild(resultItem);
            });

            // Show pagination if needed
            if (data.total_pages > 1) {
                displayPagination(data);
            } else {
                document.getElementById('pagination').style.display = 'none';
            }

            // Re-initialize lazy loading for search results
            reinitLazyLoading();
        }

        function createResultItem(result) {
            const item = document.createElement('div');
            item.className = 'post';
            item.onclick = () => window.location.href = '/post/' + result.post.id;

            const image = result.post.image || 'https://placehold.co/600x400?text=' + encodeURIComponent(result.post.title);
            const score = result.score ? 'Score: ' + result.score.toFixed(2) : '';

            item.innerHTML = 
                '<div class="image-container">' +
                '<div class="image-placeholder"></div>' +
                '<img class="lazy-image" data-src="' + image + '" alt="' + result.post.title + '" onerror="this.src=\'https://placehold.co/600x400?text=' + encodeURIComponent(result.post.title) + '\'">' +
                '</div>' +
                '<h2>' + result.post.title + '</h2>' +
                '<div class="excerpt">' + (result.post.excerpt || 'No excerpt available') + '</div>' +
                '<div class="meta">Created: ' + new Date(result.post.created_at).toLocaleDateString() + '</div>' +
                '<div class="read-more">' + score + '</div>';

            return item;
        }

        function displayPagination(data) {
            const pagination = document.getElementById('pagination');
            pagination.innerHTML = '';
            pagination.style.display = 'flex';

            // Previous button
            const prevButton = document.createElement('button');
            prevButton.textContent = '← Previous';
            prevButton.disabled = data.page <= 1;
            prevButton.onclick = () => {
                if (data.page > 1) {
                    currentPage = data.page - 1;
                    performSearch();
                }
            };
            pagination.appendChild(prevButton);

            // Page numbers
            const startPage = Math.max(1, data.page - 2);
            const endPage = Math.min(data.total_pages, data.page + 2);

            for (let i = startPage; i <= endPage; i++) {
                const pageButton = document.createElement('button');
                pageButton.textContent = i;
                pageButton.className = i === data.page ? 'current-page' : '';
                pageButton.onclick = () => {
                    currentPage = i;
                    performSearch();
                };
                pagination.appendChild(pageButton);
            }

            // Next button
            const nextButton = document.createElement('button');
            nextButton.textContent = 'Next →';
            nextButton.disabled = data.page >= data.total_pages;
            nextButton.onclick = () => {
                if (data.page < data.total_pages) {
                    currentPage = data.page + 1;
                    performSearch();
                }
            };
            pagination.appendChild(nextButton);
        }

        function showAllPosts() {
            isSearchMode = false;
            currentQuery = '';
            currentPage = 1;
            
            document.getElementById('allPosts').style.display = 'grid';
            document.getElementById('searchResults').style.display = 'none';
            document.getElementById('resultsInfo').style.display = 'none';
            document.getElementById('noResults').style.display = 'none';
            document.getElementById('pagination').style.display = 'none';
            document.getElementById('loading').style.display = 'none';
            document.getElementById('searchButton').disabled = false;
            
            document.getElementById('allPostsBtn').classList.add('active');
            document.getElementById('searchResultsBtn').classList.remove('active');
            document.getElementById('searchResultsBtn').style.display = 'none';
        }

        function showSearchResults() {
            document.getElementById('allPosts').style.display = 'none';
            document.getElementById('searchResults').style.display = 'grid';
            
            document.getElementById('allPostsBtn').classList.remove('active');
            document.getElementById('searchResultsBtn').classList.add('active');
            document.getElementById('searchResultsBtn').style.display = 'inline-block';
        }

        function showLoading() {
            document.getElementById('loading').style.display = 'block';
            document.getElementById('searchButton').disabled = true;
        }

        function hideLoading() {
            document.getElementById('loading').style.display = 'none';
            document.getElementById('searchButton').disabled = false;
        }

        function showError(message) {
            const noResults = document.getElementById('noResults');
            noResults.innerHTML = '<h3>Error</h3><p>' + message + '</p>';
            noResults.style.display = 'block';
        }

        // Lazy loading functionality
        function initLazyLoading() {
            const lazyImages = document.querySelectorAll('.lazy-image');
            console.log('Initializing lazy loading for', lazyImages.length, 'images');
            
            if ('IntersectionObserver' in window) {
                const imageObserver = new IntersectionObserver((entries, observer) => {
                    entries.forEach(entry => {
                        if (entry.isIntersecting) {
                            const img = entry.target;
                            const placeholder = img.parentElement.querySelector('.image-placeholder');
                            console.log('Image intersecting:', img.dataset.src);
                            
                            img.onload = function() {
                                console.log('Image loaded successfully:', img.src);
                                img.classList.add('loaded');
                                if (placeholder) {
                                    placeholder.style.display = 'none';
                                }
                            };
                            
                            img.onerror = function() {
                                console.log('Image failed to load:', img.dataset.src);
                                if (placeholder) {
                                    placeholder.style.display = 'none';
                                }
                            };
                            
                            img.src = img.dataset.src;
                            observer.unobserve(img);
                        }
                    });
                });

                lazyImages.forEach(img => {
                    console.log('Observing image:', img.dataset.src);
                    imageObserver.observe(img);
                });
            } else {
                // Fallback for browsers without IntersectionObserver
                console.log('IntersectionObserver not supported, loading all images immediately');
                lazyImages.forEach(img => {
                    const placeholder = img.parentElement.querySelector('.image-placeholder');
                    
                    img.onload = function() {
                        img.classList.add('loaded');
                        if (placeholder) {
                            placeholder.style.display = 'none';
                        }
                    };
                    
                    img.onerror = function() {
                        if (placeholder) {
                            placeholder.style.display = 'none';
                        }
                    };
                    
                    img.src = img.dataset.src;
                });
            }
        }

        // Initialize lazy loading when DOM is loaded
        document.addEventListener('DOMContentLoaded', function() {
            initLazyLoading();
        });

        // Re-initialize lazy loading after search results are displayed
        function reinitLazyLoading() {
            setTimeout(initLazyLoading, 100);
        }
    </script>
</body>
</html>`

	return html
}

// generatePostDetailHTML generates the post detail page HTML
func generatePostDetailHTML(post *domain.Post) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - Blog</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; }
        .back-link { margin-bottom: 20px; }
        .back-link a { color: #666; text-decoration: none; }
        .back-link a:hover { color: #333; }
        .post { background: white; border-radius: 10px; padding: 30px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .post img { width: 100%%; max-height: 400px; object-fit: cover; border-radius: 10px; margin-bottom: 20px; }
        .post .image-container { position: relative; width: 100%%; max-height: 400px; margin-bottom: 20px; border-radius: 10px; background-color: #f0f0f0; overflow: hidden; }
        .post .image-placeholder { width: 100%%; height: 100%%; background: linear-gradient(90deg, #f0f0f0 25%%, #e0e0e0 50%%, #f0f0f0 75%%); background-size: 200%% 100%%; animation: loading 1.5s infinite; }
        .post .lazy-image { width: 100%%; height: 100%%; object-fit: cover; opacity: 0; transition: opacity 0.3s ease-in-out; }
        .post .lazy-image.loaded { opacity: 1; }
        @keyframes loading {
            0%% { background-position: 200%% 0; }
            100%% { background-position: -200%% 0; }
        }
        .post h1 { color: #333; margin: 0 0 15px 0; font-size: 2.2em; line-height: 1.2; }
        .post .meta { color: #999; font-size: 0.9em; margin-bottom: 20px; padding-bottom: 20px; border-bottom: 1px solid #eee; }
        .post .content { color: #333; line-height: 1.8; font-size: 1.1em; }
        .post .content p { margin-bottom: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="back-link">
            <a href="/">← Back to Blog</a>
        </div>
        
        <div class="post">
            <div class="image-container">
                <div class="image-placeholder"></div>
                <img class="lazy-image" data-src="%s" alt="%s" onerror="this.src='https://placehold.co/600x400?text=%s'">
            </div>
            <h1>%s</h1>
            <div class="meta">Published on %s</div>
            <div class="content">
                %s
            </div>
        </div>
    </div>

    <script>
        // Lazy loading functionality for post detail page
        document.addEventListener('DOMContentLoaded', function() {
            const lazyImage = document.querySelector('.lazy-image');
            console.log('Post detail: Found lazy image:', lazyImage);
            
            if (lazyImage) {
                const placeholder = lazyImage.parentElement.querySelector('.image-placeholder');
                console.log('Post detail: Found placeholder:', placeholder);
                
                if ('IntersectionObserver' in window) {
                    console.log('Post detail: Using IntersectionObserver');
                    const imageObserver = new IntersectionObserver((entries) => {
                        entries.forEach(entry => {
                            if (entry.isIntersecting) {
                                const img = entry.target;
                                console.log('Post detail: Image intersecting:', img.dataset.src);
                                
                                img.onload = function() {
                                    console.log('Post detail: Image loaded successfully:', img.src);
                                    img.classList.add('loaded');
                                    if (placeholder) {
                                        placeholder.style.display = 'none';
                                    }
                                };
                                
                                img.onerror = function() {
                                    console.log('Post detail: Image failed to load:', img.dataset.src);
                                    if (placeholder) {
                                        placeholder.style.display = 'none';
                                    }
                                };
                                
                                img.src = img.dataset.src;
                                imageObserver.unobserve(img);
                            }
                        });
                    });
                    imageObserver.observe(lazyImage);
                } else {
                    // Fallback for browsers without IntersectionObserver
                    console.log('Post detail: IntersectionObserver not supported, loading immediately');
                    lazyImage.onload = function() {
                        console.log('Post detail: Image loaded successfully (fallback):', lazyImage.src);
                        lazyImage.classList.add('loaded');
                        if (placeholder) {
                            placeholder.style.display = 'none';
                        }
                    };
                    
                    lazyImage.onerror = function() {
                        console.log('Post detail: Image failed to load (fallback):', lazyImage.dataset.src);
                        if (placeholder) {
                            placeholder.style.display = 'none';
                        }
                    };
                    
                    lazyImage.src = lazyImage.dataset.src;
                }
            }
        });
    </script>
</body>
</html>`, post.Title, post.Image, post.Title, post.Title, post.Title, post.CreatedAt.Format("January 02, 2006"), post.Body)
}

// generateDashboardHTML generates the admin dashboard HTML
func generateDashboardHTML(posts []*domain.Post) string {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Blog Admin Dashboard</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; }
        .header { text-align: center; margin-bottom: 40px; }
        .header h1 { color: #333; font-size: 2.5em; margin-bottom: 10px; }
        .header p { color: #666; font-size: 1.2em; }
        .actions { text-align: center; margin-bottom: 30px; }
        .btn { display: inline-block; padding: 12px 24px; background-color: #007bff; color: white; text-decoration: none; border-radius: 5px; margin: 0 10px; }
        .btn:hover { background-color: #0056b3; }
        .back-link { text-align: center; margin-bottom: 20px; }
        .back-link a { color: #666; text-decoration: none; }
        .back-link a:hover { color: #333; }
        .posts { display: grid; grid-template-columns: repeat(auto-fill, minmax(350px, 1fr)); gap: 20px; }
        .post { background: white; border-radius: 10px; padding: 20px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .post img { width: 100%; height: 200px; object-fit: cover; border-radius: 5px; margin-bottom: 15px; }
        .post h2 { color: #333; margin: 0 0 10px 0; font-size: 1.4em; }
        .post .excerpt { color: #666; margin-bottom: 15px; line-height: 1.5; }
        .post .meta { color: #999; font-size: 0.9em; margin-bottom: 15px; }
        .post .actions { text-align: right; }
        .post .btn-small { padding: 8px 16px; font-size: 0.9em; }
        .btn-danger { background-color: #dc3545; }
        .btn-danger:hover { background-color: #c82333; }
        .btn-warning { background-color: #ffc107; color: #212529; }
        .btn-warning:hover { background-color: #e0a800; }
        .post .image-container { position: relative; width: 100%; height: 200px; margin-bottom: 15px; border-radius: 5px; background-color: #f0f0f0; overflow: hidden; }
        .post .image-placeholder { width: 100%; height: 100%; background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%); background-size: 200% 100%; animation: loading 1.5s infinite; }
        .post .lazy-image { width: 100%; height: 100%; object-fit: cover; opacity: 0; transition: opacity 0.3s ease-in-out; }
        .post .lazy-image.loaded { opacity: 1; }
        @keyframes loading {
            0% { background-position: 200% 0; }
            100% { background-position: -200% 0; }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Blog Admin Dashboard</h1>
            <p>Manage your blog posts</p>
        </div>
        
        <div class="back-link">
            <a href="/">← Back to Blog</a>
        </div>
        
        <div class="actions">
            <a href="/dashboard/create" class="btn">Create New Post</a>
        </div>
        
        <div class="posts">
`

	if len(posts) == 0 {
		html += `
            <div style="grid-column: 1 / -1; text-align: center; padding: 40px; color: #666;">
                <h3>No posts yet</h3>
                <p>Create your first blog post to get started!</p>
            </div>
        `
	} else {
		for _, post := range posts {
			html += fmt.Sprintf(`
            <div class="post">
                <div class="image-container">
                    <div class="image-placeholder"></div>
                    <img class="lazy-image" data-src="%s" alt="%s" onerror="this.src='https://placehold.co/600x400?text=%s'">
                </div>
                <h2>%s</h2>
                <div class="excerpt">%s</div>
                <div class="meta">Created: %s</div>
                <div class="actions">
                    <a href="/dashboard/edit?id=%d" class="btn btn-small btn-warning">Edit</a>
                    <button onclick="deletePost(%d)" class="btn btn-small btn-danger">Delete</button>
                </div>
            </div>
            `, post.Image, post.Title, post.Title, post.Title, post.Excerpt, post.CreatedAt.Format("Jan 02, 2006"), post.ID, post.ID)
		}
	}

	html += `
        </div>
    </div>

    <script>
        function deletePost(id) {
            if (confirm('Are you sure you want to delete this post?')) {
                fetch('/api/posts?id=' + id, { method: 'DELETE' })
                    .then(response => {
                        if (response.ok) {
                            location.reload();
                        } else {
                            alert('Failed to delete post');
                        }
                    })
                    .catch(error => {
                        console.error('Error:', error);
                        alert('Failed to delete post');
                    });
            }
        }

        // Lazy loading functionality for dashboard
        document.addEventListener('DOMContentLoaded', function() {
            const lazyImages = document.querySelectorAll('.lazy-image');
            console.log('Dashboard: Initializing lazy loading for', lazyImages.length, 'images');
            
            if ('IntersectionObserver' in window) {
                const imageObserver = new IntersectionObserver((entries, observer) => {
                    entries.forEach(entry => {
                        if (entry.isIntersecting) {
                            const img = entry.target;
                            const placeholder = img.parentElement.querySelector('.image-placeholder');
                            console.log('Dashboard: Image intersecting:', img.dataset.src);
                            
                            img.onload = function() {
                                console.log('Dashboard: Image loaded successfully:', img.src);
                                img.classList.add('loaded');
                                if (placeholder) {
                                    placeholder.style.display = 'none';
                                }
                            };
                            
                            img.onerror = function() {
                                console.log('Dashboard: Image failed to load:', img.dataset.src);
                                if (placeholder) {
                                    placeholder.style.display = 'none';
                                }
                            };
                            
                            img.src = img.dataset.src;
                            observer.unobserve(img);
                        }
                    });
                });

                lazyImages.forEach(img => {
                    console.log('Dashboard: Observing image:', img.dataset.src);
                    imageObserver.observe(img);
                });
            } else {
                // Fallback for browsers without IntersectionObserver
                console.log('Dashboard: IntersectionObserver not supported, loading all images immediately');
                lazyImages.forEach(img => {
                    const placeholder = img.parentElement.querySelector('.image-placeholder');
                    
                    img.onload = function() {
                        img.classList.add('loaded');
                        if (placeholder) {
                            placeholder.style.display = 'none';
                        }
                    };
                    
                    img.onerror = function() {
                        if (placeholder) {
                            placeholder.style.display = 'none';
                        }
                    };
                    
                    img.src = img.dataset.src;
                });
            }
        });
    </script>
</body>
</html>`

	return html
}

// generateCreateFormHTML generates the create post form HTML
func generateCreateFormHTML() string {
	return `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Create New Post</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; }
        .header { text-align: center; margin-bottom: 30px; }
        .header h1 { color: #333; }
        .form-group { margin-bottom: 20px; }
        label { display: block; margin-bottom: 5px; color: #333; font-weight: bold; }
        input[type="text"], textarea { width: 100%; padding: 10px; border: 1px solid #ddd; border-radius: 5px; font-size: 16px; }
        textarea { height: 120px; resize: vertical; }
        .btn { padding: 12px 24px; background-color: #007bff; color: white; border: none; border-radius: 5px; cursor: pointer; font-size: 16px; }
        .btn:hover { background-color: #0056b3; }
        .btn-secondary { background-color: #6c757d; margin-right: 10px; }
        .btn-secondary:hover { background-color: #545b62; }
        .actions { text-align: center; margin-top: 30px; }
        .back-link { text-align: center; margin-bottom: 20px; }
        .back-link a { color: #666; text-decoration: none; }
        .back-link a:hover { color: #333; }
    </style>
</head>
<body>
    <div class="container">
        <div class="back-link">
            <a href="/dashboard">← Back to Dashboard</a>
        </div>
        
        <div class="header">
            <h1>Create New Post</h1>
        </div>
        
        <form id="createForm">
            <div class="form-group">
                <label for="title">Title *</label>
                <input type="text" id="title" name="title" required>
            </div>
            
            <div class="form-group">
                <label for="image">Image URL</label>
                <input type="text" id="image" name="image" placeholder="https://example.com/image.jpg">
            </div>
            
            <div class="form-group">
                <label for="excerpt">Excerpt</label>
                <textarea id="excerpt" name="excerpt" placeholder="Brief summary of your post..."></textarea>
            </div>
            
            <div class="form-group">
                <label for="body">Body *</label>
                <textarea id="body" name="body" required placeholder="Write your post content here..."></textarea>
            </div>
            
            <div class="actions">
                <a href="/dashboard" class="btn btn-secondary">Cancel</a>
                <button type="submit" class="btn">Create Post</button>
            </div>
        </form>
    </div>

    <script>
        document.getElementById('createForm').addEventListener('submit', function(e) {
            e.preventDefault();
            
            const formData = {
                title: document.getElementById('title').value,
                image: document.getElementById('image').value,
                excerpt: document.getElementById('excerpt').value,
                body: document.getElementById('body').value
            };
            
            fetch('/api/posts', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(formData)
            })
            .then(response => {
                if (response.ok) {
                    alert('Post created successfully!');
                    window.location.href = '/dashboard';
                } else {
                    alert('Failed to create post');
                }
            })
            .catch(error => {
                console.error('Error:', error);
                alert('Failed to create post');
            });
        });
    </script>
</body>
</html>`
}

// generateEditFormHTML generates the edit post form HTML
func generateEditFormHTML(post *domain.Post) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Edit Post</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; }
        .header { text-align: center; margin-bottom: 30px; }
        .header h1 { color: #333; }
        .form-group { margin-bottom: 20px; }
        label { display: block; margin-bottom: 5px; color: #333; font-weight: bold; }
        input[type="text"], textarea { width: 100%%; padding: 10px; border: 1px solid #ddd; border-radius: 5px; font-size: 16px; }
        textarea { height: 120px; resize: vertical; }
        .btn { padding: 12px 24px; background-color: #007bff; color: white; border: none; border-radius: 5px; cursor: pointer; font-size: 16px; }
        .btn:hover { background-color: #0056b3; }
        .btn-secondary { background-color: #6c757d; margin-right: 10px; }
        .btn-secondary:hover { background-color: #545b62; }
        .actions { text-align: center; margin-top: 30px; }
        .back-link { text-align: center; margin-bottom: 20px; }
        .back-link a { color: #666; text-decoration: none; }
        .back-link a:hover { color: #333; }
    </style>
</head>
<body>
    <div class="container">
        <div class="back-link">
            <a href="/dashboard">← Back to Dashboard</a>
        </div>
        
        <div class="header">
            <h1>Edit Post</h1>
        </div>
        
        <form id="editForm">
            <div class="form-group">
                <label for="title">Title *</label>
                <input type="text" id="title" name="title" value="%s" required>
            </div>
            
            <div class="form-group">
                <label for="image">Image URL</label>
                <input type="text" id="image" name="image" value="%s" placeholder="https://example.com/image.jpg">
            </div>
            
            <div class="form-group">
                <label for="excerpt">Excerpt</label>
                <textarea id="excerpt" name="excerpt" placeholder="Brief summary of your post...">%s</textarea>
            </div>
            
            <div class="form-group">
                <label for="body">Body *</label>
                <textarea id="body" name="body" required placeholder="Write your post content here...">%s</textarea>
            </div>
            
            <div class="actions">
                <a href="/dashboard" class="btn btn-secondary">Cancel</a>
                <button type="submit" class="btn">Update Post</button>
            </div>
        </form>
    </div>

    <script>
        document.getElementById('editForm').addEventListener('submit', function(e) {
            e.preventDefault();
            
            const formData = {
                title: document.getElementById('title').value,
                image: document.getElementById('image').value,
                excerpt: document.getElementById('excerpt').value,
                body: document.getElementById('body').value
            };
            
            fetch('/api/posts?id=%d', {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(formData)
            })
            .then(response => {
                if (response.ok) {
                    alert('Post updated successfully!');
                    window.location.href = '/dashboard';
                } else {
                    alert('Failed to update post');
                }
            })
            .catch(error => {
                console.error('Error:', error);
                alert('Failed to update post');
            });
        });
    </script>
</body>
</html>`, post.Title, post.Image, post.Excerpt, post.Body, post.ID)
}
