# Use a pre-built Fluentd image with common plugins
FROM fluent/fluentd:v1.16-1

# Switch to root to install only essential plugins
USER root

# Install only the essential plugins we need
RUN gem install --no-document \
    fluent-plugin-elasticsearch \
    fluent-plugin-rewrite-tag-filter

# Switch back to fluent user for security
USER fluent 