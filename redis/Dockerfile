FROM redis:6.0-alpine
# RUN apk --no-cache add sudo && \
# 	adduser -D user && \
# 	echo "ALL ALL = (ALL) NOPASSWD: ALL" > /etc/sudoers
 
WORKDIR /redis
#CMD ["redis-server", "--bind", "redis", "--port", "6379"]
COPY redis.conf /usr/local/etc/redis/redis.conf
# COPY init.sh ./
# RUN sudo chmod +x init.sh
CMD [ "redis-server", "/usr/local/etc/redis/redis.conf"]