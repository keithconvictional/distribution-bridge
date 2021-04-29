if [ -z $RENDER_WEBHOOK_URL ]
then
  echo "Render Webhook not set, skipping."
  exit 0
fi

curl $RENDER_WEBHOOK_URL