import path from 'node:path'
import fs from 'node:fs'
import process from 'node:process'
import { globby } from 'globby'
import { PutObjectCommand, S3Client } from '@aws-sdk/client-s3'
import { fileTypeFromBuffer } from 'file-type'
import mime from 'mime'

export default function myPlugin(targets: { src: string; keyDir: string }[] = []) {
  const env = process.env

  const S3Config = {
    region: 'auto',
    endpoint: `https://${env.CLOUDFLARE_ACCOUNT_ID}.r2.cloudflarestorage.com`,
    credentials: {
      accessKeyId: env.CLOUDFLARE_ACCESSKEY_ID,
      secretAccessKey: env.CLOUDFLARE_ACCESSKEY_SECRET,
    },
  }

  const client = new S3Client(S3Config)

  return {
    name: 'vite-cloudflare-upload-plugin',
    writeBundle: async () => {
      if (!env.CLOUDFLARE_ACCESSKEY_ID || !env.CLOUDFLARE_ACCESSKEY_SECRET || !env.CLOUDFLARE_ACCOUNT_ID || !env.CLOUDFLARE_BUCKET_NAME) {
        console.error('Cloudflare credentials not found, skipping upload')
        return
      }

      await new Promise(res => setTimeout(res, 100))

      if (Array.isArray(targets) && targets.length) {
        const commandLog = {
          fatal: [],
          success: [],
        }
        const reqs = []
        for (const target of targets) {
          const matchedPaths = await globby(target.src, {
            expandDirectories: false,
            onlyFiles: false,
          })
          for (let i = 0; i < matchedPaths.length; i++) {
            const p = matchedPaths[i]
            const key = path.join(target.keyDir, path.basename(p))
            const file = fs.readFileSync(p)
            let type = 'application/octet-stream'
            const t1 = mime.getType(p)
            if (t1) {
              type = t1
            }
            else {
              const t2 = await fileTypeFromBuffer(file)
              if (t2)
                type = t2.mime
            }

            const command = new PutObjectCommand({
              Bucket: env.CLOUDFLARE_BUCKET_NAME,
              Key: key,
              Body: file,
              ContentType: type,
            })
            reqs.push(
              client.send(command).then((res) => {
                commandLog.success.push({ p, key, res })
              }).catch((e) => {
                commandLog.fatal.push({ p, key, e })
              }),
            )
          }
        }
        await Promise.all(reqs)
        fs.writeFileSync('./commandLog.json', JSON.stringify(commandLog, null, 2))
      }
    },
  }
}
