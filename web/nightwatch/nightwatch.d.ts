import 'nightwatch'

declare module 'nightwatch' {
  // These interfaces are intentionally empty: they are declaration-merging
  // augmentation points for project-specific assertions/commands.
  /* eslint-disable @typescript-eslint/no-empty-object-type */
  interface NightwatchCustomAssertions {
    // Add your custom assertions' types here
    // elementHasCount: (selector: string, count: number) => NightwatchBrowser
  }

  interface NightwatchCustomCommands {
    // Add your custom commands' types here
    // strictClick: (selector: string) => NightwatchBrowser
  }
  /* eslint-enable @typescript-eslint/no-empty-object-type */
}
