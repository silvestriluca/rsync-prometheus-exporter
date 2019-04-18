#!/bin/bash
rsync --daemon -v --config rsyncd_example.conf
rsync -rav rsync://localhost:8737/exercise/ tmp/
app