type AppErrorParams = {
  status: number
  message: string
  description?: string
}

class AppError extends Error {
  status: number
  description?: string

  constructor({ status, message, description }: AppErrorParams) {
    super(message)

    this.status = status
    this.description = description
    this.name = "AppError"

    Object.setPrototypeOf(this, new.target.prototype)
  }
}

export { AppError, type AppErrorParams }
